package bank

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tokenize-titan/titan/utils"

	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/testutil/cmd/tx"
	txcmd "github.com/tokenize-titan/titan/testutil/cmd/tx"
)

func MustSend(t testing.TB, from string, to string, amount string) txcmd.TxResponse {
	fromBalBefore := MustGetBalance(t, from, utils.BaseDenom, 0)
	toBalBefore := MustGetBalance(t, to, utils.BaseDenom, 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "bank", "send", from, to, amount)

	fromBalAfter := MustGetBalance(t, from, utils.BaseDenom, 0)
	toBalAfter := MustGetBalance(t, to, utils.BaseDenom, 0)

	coinSpent, err := tx.GetDeductFeeAmount()
	require.NoError(t, err)
	sentAmount := testutil.MustGetBaseDenomAmount(t, amount)

	require.Equal(t, fromBalBefore.Sub(coinSpent).Sub(sentAmount), fromBalAfter)
	require.Equal(t, toBalBefore.Add(sentAmount), toBalAfter)

	return tx
}

func MustMultiSend(t testing.TB, from string, amount string, to ...string) tx.TxResponse {
	fromBalBefore := MustGetBalance(t, from, utils.BaseDenom, 0)
	var toBalBefore []testutil.BigInt
	for i := range to {
		toBalBefore = append(toBalBefore, MustGetBalance(t, to[i], utils.BaseDenom, 0))
	}

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	args := []string{"bank", "multi-send", from}
	args = append(args, to...)
	args = append(args, amount)
	tx := txcmd.MustExecTx(t, ctx, args...)

	fromBalAfter := MustGetBalance(t, from, utils.BaseDenom, 0)
	var toBalAfter []testutil.BigInt
	for i := range to {
		toBalAfter = append(toBalAfter, MustGetBalance(t, to[i], utils.BaseDenom, 0))
	}

	coinSpent, err := tx.GetDeductFeeAmount()
	require.NoError(t, err)
	sentAmount := testutil.MustGetBaseDenomAmount(t, amount)
	totalSentAmount := sentAmount.Mul(testutil.MakeBigInt(int64(len(to))))

	require.Equal(t, fromBalBefore.Sub(coinSpent).Sub(totalSentAmount), fromBalAfter)
	for i := range to {
		require.Equal(t, toBalBefore[i].Add(sentAmount), toBalAfter[i])
	}

	return tx
}

func MustGetBalance(t testing.TB, address string, denom string, height int64) testutil.BigInt {
	args := []string{
		"bank",
		"balances",
		address,
		"--denom=" + denom,
	}
	if height > 0 {
		args = append(args, "--height="+testutil.FormatInt(height))
	}
	var data struct {
		Amount testutil.BigInt `json:"amount"`
	}
	cmd.MustQuery(t, &data, args...)
	return data.Amount
}
