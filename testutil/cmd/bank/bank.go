package bank

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/testutil/cmd/tx"
	txcmd "github.com/tokenize-titan/titan/testutil/cmd/tx"
)

func MustSend(t testing.TB, from string, to string, amount string) tx.Tx {
	fromBalBefore := MustGetBalance(t, from, "utkx")
	toBalBefore := MustGetBalance(t, to, "utkx")

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "bank", "send", from, to, amount)

	fromBalAfter := MustGetBalance(t, from, "utkx")
	toBalAfter := MustGetBalance(t, to, "utkx")

	coinSpent := tx.GasWanted.Mul(testutil.MakeBigInt(10)) // Gas price == 10 utkx
	sentAmount := testutil.MustGetUtkxAmount(t, amount)

	require.Equal(t, fromBalBefore.Sub(coinSpent).Sub(sentAmount), fromBalAfter)
	require.Equal(t, toBalBefore.Add(sentAmount), toBalAfter)

	return tx
}

func MustGetBalance(t testing.TB, address string, denom string) testutil.BigInt {
	var data struct {
		Amount testutil.BigInt `json:"amount"`
	}
	cmd.MustQuery(t, &data, "bank", "balances", address, "--denom="+denom)
	return data.Amount
}
