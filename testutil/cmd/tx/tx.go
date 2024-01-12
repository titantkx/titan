package tx

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/testutil/cmd/feemarket"
	"github.com/tokenize-titan/titan/utils"
)

var rpcErrPattern = regexp.MustCompile(`RPC\serror\s(-?[\d]+)`)

type TxResponse struct {
	Height    testutil.Int    `json:"height"`
	Code      int             `json:"code"`
	Hash      string          `json:"txhash"`
	RawLog    string          `json:"raw_log"`
	GasUsed   testutil.BigInt `json:"gas_used"`
	GasWanted testutil.BigInt `json:"gas_wanted"`
	Tx        Tx              `json:"tx"`
	Events    []Event         `json:"events"`
}

type Tx struct {
	Type       string   `json:"@type"`
	Body       struct{} `json:"body"`
	AuthInfo   AuthInfo `json:"auth_info"`
	Signatures []string `json:"signatures"`
}

type AuthInfo struct {
	SignerInfos []struct{} `json:"signer_infos"`
	Fee         Fee        `json:"fee"`
	Tip         *struct{}  `json:"tip"`
}

type Fee struct {
	Amount   testutil.Coins  `json:"amount"`
	GasLimit testutil.BigInt `json:"gas_limit"`
	Payer    string          `json:"payer"`
	Granter  string          `json:"granter"`
}

type Event struct {
	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Index bool   `json:"index"`
}

func QueryTx(ctx context.Context, txHash string) (*TxResponse, error) {
	for {
		output, err := cmd.Exec("titand", "query", "tx", txHash, "--output=json")
		if err != nil {
			matches := rpcErrPattern.FindStringSubmatch(string(output))
			if len(matches) == 2 && matches[1] == "-32603" {
				// Transaction not found, wait until it is delivered or timeout
				select {
				case <-time.After(1 * time.Second):
					continue
				case <-ctx.Done():
					if err == nil {
						err = ctx.Err()
					}
					return nil, err
				}
			}
			return nil, err
		}
		var tx TxResponse
		if err := cmd.UnmarshalJSON(output, &tx); err != nil {
			return nil, err
		}
		return &tx, nil
	}
}

func ExecTx(ctx context.Context, args ...string) (*TxResponse, error) {
	gasPrice, err := feemarket.GetBaseFee(0)
	if err != nil {
		return nil, err
	}
	args = append([]string{"tx"}, args...)
	args = append(args, "--gas=auto", "--gas-adjustment=1.3", "--gas-prices="+gasPrice.String()+utils.BaseDenom, "--output=json", "-y")
	args = append(args, "--keyring-backend=test")
	output, err := cmd.Exec("titand", args...)
	if err != nil {
		return nil, err
	}
	var tx TxResponse
	if err := cmd.UnmarshalJSON(output, &tx); err != nil {
		return nil, err
	}
	if tx.Code != 0 {
		return nil, errors.New(tx.RawLog)
	}
	return QueryTx(ctx, tx.Hash)
}

func MustExecTx(t testing.TB, ctx context.Context, args ...string) TxResponse {
	tx, err := ExecTx(ctx, args...)
	require.NoError(t, err)
	require.Equal(t, 0, tx.Code, tx.RawLog)
	return *tx
}

func (txr TxResponse) GetRefundAmount() (testutil.BigInt, error) {
	// find event in `tx.Events` have type "refund"
	var refundEvent *Event
	for _, event := range txr.Events {
		if event.Type == "refund" {
			refundEvent = &event
			break
		}
	}

	if refundEvent != nil {
		// find attribute "amount" in `refundEvent.Attributes`
		var amountValue string
		for _, attr := range refundEvent.Attributes {
			if attr.Key == "amount" {
				amountValue = attr.Value
				break
			}
		}

		// convert amount value to BigInt
		if amountValue != "" {
			refundAmount, err := testutil.ParseAmount(amountValue)
			if err == nil {
				return refundAmount.GetBaseDenomAmount(), err
			}
		}

		return testutil.Coins{}.GetBaseDenomAmount(), nil
	}

	return testutil.Coins{}.GetBaseDenomAmount(), nil
}

func (txr TxResponse) MustGetRefundAmount(t testing.TB) testutil.BigInt {
	amount, err := txr.GetRefundAmount()
	require.NoError(t, err)
	return amount
}

func (txr TxResponse) GetDeductFeeAmount() (testutil.BigInt, error) {
	coinSpent := txr.Tx.AuthInfo.Fee.Amount.GetBaseDenomAmount()

	refundAmount, err := txr.GetRefundAmount()
	if err != nil {
		return testutil.Coins{}.GetBaseDenomAmount(), err
	}

	return coinSpent.Sub(refundAmount), nil
}

func (txr TxResponse) MustGetDeductFeeAmount(t testing.TB) testutil.BigInt {
	amount, err := txr.GetDeductFeeAmount()
	require.NoError(t, err)
	return amount
}

func GenerateTx(args ...string) ([]byte, error) {
	gasPrice, err := feemarket.GetBaseFee(0)
	if err != nil {
		return nil, err
	}
	args = append([]string{"tx"}, args...)
	args = append(args, "--gas=auto", "--gas-adjustment=2", "--gas-prices="+gasPrice.String()+utils.BaseDenom, "--generate-only")
	output, err := cmd.Exec("titand", args...)
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(bytes.NewBuffer(output))
	if s.Scan() && s.Scan() {
		return s.Bytes(), nil
	}
	return nil, fmt.Errorf("invalid output: %s", string(output))
}

func MustGenerateTx(t testing.TB, args ...string) []byte {
	output, err := GenerateTx(args...)
	require.NoError(t, err)
	require.NotEmpty(t, output)
	return output
}

func BroadcastTx(ctx context.Context, filePath string) (*TxResponse, error) {
	output, err := cmd.Exec("titand", "tx", "broadcast", filePath, "--output=json", "-y")
	if err != nil {
		return nil, err
	}
	var tx TxResponse
	if err := cmd.UnmarshalJSON(output, &tx); err != nil {
		return nil, err
	}
	if tx.Code != 0 {
		return nil, errors.New(tx.RawLog)
	}
	return QueryTx(ctx, tx.Hash)
}

func MustBroadcastTx(t testing.TB, ctx context.Context, filePath string) TxResponse {
	tx, err := BroadcastTx(ctx, filePath)
	require.NoError(t, err)
	require.Equal(t, 0, tx.Code, tx.RawLog)
	return *tx
}
