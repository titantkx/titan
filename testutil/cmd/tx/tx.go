package tx

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/testutil/cmd/feemarket"
	"github.com/tokenize-titan/titan/utils"
)

var rpcErrPattern = regexp.MustCompile(`RPC\serror\s(-?[\d]+)`)

type TxResponse struct {
	Height    testutil.Int `json:"height"`
	Code      int          `json:"code"`
	Hash      string       `json:"txhash"`
	RawLog    string       `json:"raw_log"`
	GasUsed   testutil.Int `json:"gas_used"`
	GasWanted testutil.Int `json:"gas_wanted"`
	Tx        Tx           `json:"tx"`
	Events    []Event      `json:"events"`
}

func (txr TxResponse) FindEvent(typ string) *Event {
	for _, evt := range txr.Events {
		if evt.Type == typ {
			return &evt
		}
	}
	return nil
}

func (txr TxResponse) MustGetEventAttributeValue(t testutil.TestingT, eventType string, attributeKey string) string {
	event := txr.FindEvent(eventType)
	require.NotNil(t, event)
	attr := event.FindAttribute(attributeKey)
	require.NotNil(t, attr)
	return attr.Value
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
	Amount   testutil.Coins `json:"amount"`
	GasLimit testutil.Int   `json:"gas_limit"`
	Payer    string         `json:"payer"`
	Granter  string         `json:"granter"`
}

type Event struct {
	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

func (evt Event) FindAttribute(key string) *Attribute {
	for _, attr := range evt.Attributes {
		if attr.Key == key {
			return &attr
		}
	}
	return nil
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

func MustExecTx(t testutil.TestingT, ctx context.Context, args ...string) TxResponse {
	tx, err := ExecTx(ctx, args...)
	require.NoError(t, err)
	require.Equal(t, 0, tx.Code, tx.RawLog)
	return *tx
}

func MustErrExecTx(t testutil.TestingT, ctx context.Context, expErr string, args ...string) {
	tx, err := ExecTx(ctx, args...)
	require.Nil(t, tx)
	require.Error(t, err)
	if expErr != "" {
		require.ErrorContains(t, err, expErr)
	}
}

func (txr TxResponse) GetRefundAmount() (testutil.Int, error) {
	refundEvent := txr.FindEvent("refund")
	if refundEvent == nil {
		return testutil.MakeInt(0), nil
	}

	amountAttr := refundEvent.FindAttribute("amount")
	if amountAttr == nil {
		return testutil.MakeInt(0), errors.New("amount attribute is required")
	}

	amount, err := testutil.ParseAmount(amountAttr.Value)
	if err != nil {
		return testutil.MakeInt(0), err
	}

	return amount.GetBaseDenomAmount(), nil
}

func (txr TxResponse) MustGetRefundAmount(t testutil.TestingT) testutil.Int {
	amount, err := txr.GetRefundAmount()
	require.NoError(t, err)
	return amount
}

func (txr TxResponse) GetDeductFeeAmount() (testutil.Int, error) {
	coinSpent := txr.Tx.AuthInfo.Fee.Amount.GetBaseDenomAmount()

	refundAmount, err := txr.GetRefundAmount()
	if err != nil {
		return testutil.Coins{}.GetBaseDenomAmount(), err
	}

	return coinSpent.Sub(refundAmount), nil
}

func (txr TxResponse) MustGetDeductFeeAmount(t testutil.TestingT) testutil.Int {
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

func MustGenerateTx(t testutil.TestingT, args ...string) []byte {
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

func MustBroadcastTx(t testutil.TestingT, ctx context.Context, filePath string) TxResponse {
	tx, err := BroadcastTx(ctx, filePath)
	require.NoError(t, err)
	require.Equal(t, 0, tx.Code, tx.RawLog)
	return *tx
}
