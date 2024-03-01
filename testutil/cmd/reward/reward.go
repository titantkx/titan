package reward

import (
	"context"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/testutil/cmd/bank"
	txcmd "github.com/tokenize-titan/titan/testutil/cmd/tx"
	"github.com/tokenize-titan/titan/utils"
)

type Params struct {
	Authority string         `json:"authority"`
	Rate      testutil.Float `json:"rate"`
}

func MustGetParams(t testutil.TestingT) Params {
	var v struct {
		Params Params `json:"params"`
	}
	cmd.MustQuery(t, &v, "validatorreward", "params")
	return v.Params
}

func MustFundRewardPool(t testutil.TestingT, from string, amount string) {
	balBefore := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "validatorreward", "fund-reward-pool", amount, "--from="+from)

	balAfter := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	coinSpent := tx.MustGetDeductFeeAmount(t)
	fundAmount := testutil.MustGetBaseDenomAmount(t, amount)

	require.Equal(t, balBefore.Sub(coinSpent).Sub(fundAmount), balAfter)
}

func SetRate(from string, newRate testutil.Float) error {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()
	_, err := txcmd.ExecTx(ctx, "validatorreward", "set-rate", newRate.String(), "--from="+from)
	return err
}

func MustSetRate(t testutil.TestingT, from string, newRate testutil.Float) {
	err := SetRate(from, newRate)
	require.NoError(t, err)
	params := MustGetParams(t)
	require.Equal(t, from, params.Authority)
	require.Equal(t, newRate.String(), params.Rate.String())
}

func SetAuthority(from string, newAuthority string) error {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()
	_, err := txcmd.ExecTx(ctx, "validatorreward", "set-authority", newAuthority, "--from="+from)
	return err
}

func MustSetAuthority(t testutil.TestingT, from string, newAuthority string) {
	err := SetAuthority(from, newAuthority)
	require.NoError(t, err)
	params := MustGetParams(t)
	require.Equal(t, newAuthority, params.Authority)
}
