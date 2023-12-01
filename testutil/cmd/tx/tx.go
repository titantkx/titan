package tx

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
)

var rpcErrPattern = regexp.MustCompile("RPC\\serror\\s(-?[\\d]+)")

type Tx struct {
	Height    testutil.Int    `json:"height"`
	Code      int             `json:"code"`
	Hash      string          `json:"txhash"`
	RawLog    string          `json:"raw_log"`
	GasUsed   testutil.BigInt `json:"gas_used"`
	GasWanted testutil.BigInt `json:"gas_wanted"`
	Events    []Event         `json:"events"`
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

func ParseTx(buf []byte) (*Tx, error) {
	s := bufio.NewScanner(bytes.NewBuffer(buf))
	if !s.Scan() || !s.Scan() {
		return nil, errors.New("cannot parse Tx")
	}
	var tx Tx
	if err := json.Unmarshal(s.Bytes(), &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

func QueryTx(ctx context.Context, txHash string) (*Tx, error) {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	var err error
	for {
		select {
		case <-tick.C:
			output, err := cmd.Exec("titand", "query", "tx", txHash, "--output=json")
			if err != nil {
				matches := rpcErrPattern.FindStringSubmatch(string(output))
				if len(matches) == 2 && matches[1] == "-32603" {
					// Transaction not found, wait until it is delivered or timeout
					continue
				}
				return nil, err
			}
			var tx Tx
			if err := json.Unmarshal(output, &tx); err != nil {
				return nil, err
			}
			return &tx, nil
		case <-ctx.Done():
			if err == nil {
				err = ctx.Err()
			}
			return nil, err
		}
	}
}

func ExecTx(ctx context.Context, args ...string) (*Tx, error) {
	args = append([]string{"tx"}, args...)
	args = append(args, "--gas=auto", "--gas-adjustment=2", "--gas-prices=10utkx", "--output=json", "-y")
	output, err := cmd.Exec("titand", args...)
	if err != nil {
		return nil, err
	}
	tx, err := ParseTx(output)
	if err != nil {
		return nil, err
	}
	return QueryTx(ctx, tx.Hash)
}

func MustExecTx(t testing.TB, ctx context.Context, args ...string) Tx {
	tx, err := ExecTx(ctx, args...)
	require.NoError(t, err)
	require.Equal(t, 0, tx.Code)
	return *tx
}
