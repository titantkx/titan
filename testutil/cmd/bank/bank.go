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

func MustMultiSend(t testing.TB, from string, amount string, to ...string) tx.Tx {
	fromBalBefore := MustGetBalance(t, from, "utkx")
	var toBalBefore []testutil.BigInt
	for i := range to {
		toBalBefore = append(toBalBefore, MustGetBalance(t, to[i], "utkx"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	args := []string{"bank", "multi-send", from}
	args = append(args, to...)
	args = append(args, amount)
	tx := txcmd.MustExecTx(t, ctx, args...)

	fromBalAfter := MustGetBalance(t, from, "utkx")
	var toBalAfter []testutil.BigInt
	for i := range to {
		toBalAfter = append(toBalAfter, MustGetBalance(t, to[i], "utkx"))
	}

	coinSpent := tx.GasWanted.Mul(testutil.MakeBigInt(10)) // Gas price == 10 utkx
	sentAmount := testutil.MustGetUtkxAmount(t, amount)
	totalSentAmount := sentAmount.Mul(testutil.MakeBigInt(int64(len(to))))

	require.Equal(t, fromBalBefore.Sub(coinSpent).Sub(totalSentAmount), fromBalAfter)
	for i := range to {
		require.Equal(t, toBalBefore[i].Add(sentAmount), toBalAfter[i])
	}

	return tx
}

func MustGetBalance(t testing.TB, address string, denom string) testutil.BigInt {
	var data struct {
		Amount testutil.BigInt `json:"amount"`
	}
	cmd.MustQuery(t, &data, "bank", "balances", address, "--denom="+denom)
	return data.Amount
}
