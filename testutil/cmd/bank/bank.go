package bank

import (
	"context"

	"github.com/stretchr/testify/require"

	"github.com/titantkx/titan/utils"

	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
	txcmd "github.com/titantkx/titan/testutil/cmd/tx"
)

func MustSend(t testutil.TestingT, from string, to string, amount string) txcmd.TxResponse {
	fromBalBefore := MustGetBalance(t, from, utils.BaseDenom, 0)
	toBalBefore := MustGetBalance(t, to, utils.BaseDenom, 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "bank", "send", from, to, amount)

	fromBalAfter := MustGetBalance(t, from, utils.BaseDenom, 0)
	toBalAfter := MustGetBalance(t, to, utils.BaseDenom, 0)

	coinSpent := tx.MustGetDeductFeeAmount(t)
	sentAmount := testutil.MustGetBaseDenomAmount(t, amount)

	require.Equal(t, fromBalBefore.Sub(coinSpent).Sub(sentAmount), fromBalAfter)
	require.Equal(t, toBalBefore.Add(sentAmount), toBalAfter)

	return tx
}

func MustMultiSend(t testutil.TestingT, from string, amount string, to ...string) txcmd.TxResponse {
	fromBalBefore := MustGetBalance(t, from, utils.BaseDenom, 0)
	toBalBefore := make([]testutil.Int, 0, len(to)) // Pre-allocate toBalBefore with a capacity equal to the length of 'to'
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
	toBalAfter := make([]testutil.Int, 0, len(to)) // Pre-allocate toBalAfter with a capacity equal to the length of 'to'
	for i := range to {
		toBalAfter = append(toBalAfter, MustGetBalance(t, to[i], utils.BaseDenom, 0))
	}

	coinSpent := tx.MustGetDeductFeeAmount(t)
	sentAmount := testutil.MustGetBaseDenomAmount(t, amount)
	totalSentAmount := sentAmount.Mul(testutil.MakeInt(int64(len(to))))

	require.Equal(t, fromBalBefore.Sub(coinSpent).Sub(totalSentAmount), fromBalAfter)
	for i := range to {
		require.Equal(t, toBalBefore[i].Add(sentAmount), toBalAfter[i])
	}

	return tx
}

func MustGetBalance(t testutil.TestingT, address string, denom string, height int64) testutil.Int {
	args := []string{
		"bank",
		"balances",
		address,
		"--denom=" + denom,
	}
	if height > 0 {
		args = append(args, "--height="+testutil.FormatInt(height))
	}
	var v struct {
		Amount testutil.Int `json:"amount"`
	}
	cmd.MustQuery(t, &v, args...)
	return v.Amount
}

func MustGetTotalBalance(t testutil.TestingT, denom string, height int64) testutil.Int {
	args := []string{
		"bank",
		"total",
		"--denom=" + denom,
	}
	if height > 0 {
		args = append(args, "--height="+testutil.FormatInt(height))
	}
	var v struct {
		Amount testutil.Int `json:"amount"`
	}
	cmd.MustQuery(t, &v, args...)
	return v.Amount
}
