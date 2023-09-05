package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/cometbft/cometbft/abci/types"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/tokenize-titan/titan/app"
	"github.com/tokenize-titan/titan/x/staking/keeper"
)

var PKs = simtestutil.CreateTestPubKeys(500)

func applyValidatorSetUpdates(t *testing.T, ctx sdk.Context, k *keeper.Keeper, expectedUpdatesLen int) []abci.ValidatorUpdate {
	updates, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	if expectedUpdatesLen >= 0 {
		require.Equal(t, expectedUpdatesLen, len(updates), "%v", updates)
	}
	return updates
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(appIn *app.App, ctx sdk.Context, genAddr sdk.AccAddress, numAddrs int) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := app.AddTestAddrsIncremental(appIn, ctx, genAddr, numAddrs, sdk.NewInt(10000))
	addrVals := simtestutil.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}

func createValidators(t *testing.T, ctx sdk.Context, appIn *app.App, genAddr sdk.AccAddress, powers []int64) ([]sdk.AccAddress, []sdk.ValAddress, []types.Validator) {
	addrs := app.AddTestAddrsIncremental(appIn, ctx, genAddr, 5, appIn.StakingKeeper.TokensFromConsensusPower(ctx, 300))
	valAddrs := simtestutil.ConvertAddrsToValAddrs(addrs)
	pks := simtestutil.CreateTestPubKeys(5)

	// cdc := moduletestutil.MakeTestEncodingConfig().Codec

	// appIn.StakingKeeper = keeper.NewKeeper(
	// 	cdc,
	// 	appIn.GetKey(types.StoreKey),
	// 	appIn.AccountKeeper,
	// 	appIn.BankKeeper,
	// 	authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	// )

	// applyValidatorSetUpdates(t, ctx, appIn.StakingKeeper, -1)

	currVals := appIn.StakingKeeper.GetAllValidators(ctx)
	// print all validators addr and power
	for _, val := range currVals {
		t.Logf("val addr: %s, power: %s", val.GetOperator(), val.GetStatus().String())
	}

	t.Logf("=====================================")

	val1 := testutil.NewValidator(t, valAddrs[0], pks[0])
	val2 := testutil.NewValidator(t, valAddrs[1], pks[1])
	vals := []types.Validator{val1, val2}

	appIn.StakingKeeper.SetValidator(ctx, val1)
	// appIn.StakingKeeper.SetValidator(ctx, val2)
	appIn.StakingKeeper.SetValidatorByConsAddr(ctx, val1)
	// appIn.StakingKeeper.SetValidatorByConsAddr(ctx, val2)
	appIn.StakingKeeper.SetNewValidatorByPowerIndex(ctx, val1)
	// appIn.StakingKeeper.SetNewValidatorByPowerIndex(ctx, val2)

	// call the after-creation hook
	if err := appIn.StakingKeeper.Hooks().AfterValidatorCreated(ctx, val1.GetOperator()); err != nil {
		require.NoError(t, err)
	}

	// currVals = appIn.StakingKeeper.GetAllValidators(ctx)
	// // print all validators addr
	// for _, val := range currVals {
	// 	t.Logf("val addr: %s, power: %s", val.GetOperator(), val.GetStatus().String())
	// }
	// t.Logf("=====================================")

	_, err := appIn.StakingKeeper.Delegate(ctx, addrs[0], appIn.StakingKeeper.TokensFromConsensusPower(ctx, powers[0]), types.Unbonded, val1, true)
	require.NoError(t, err)
	// _, err = appIn.StakingKeeper.Delegate(ctx, addrs[1], appIn.StakingKeeper.TokensFromConsensusPower(ctx, powers[1]), types.Unbonded, val2, true)
	// require.NoError(t, err)
	// _, err = appIn.StakingKeeper.Delegate(ctx, addrs[0], appIn.StakingKeeper.TokensFromConsensusPower(ctx, powers[2]), types.Unbonded, val2, true)
	// require.NoError(t, err)

	// applyValidatorSetUpdates(t, ctx, appIn.StakingKeeper, -1)

	currVals = appIn.StakingKeeper.GetAllValidators(ctx)
	// print all validators addr
	for _, val := range currVals {
		t.Logf("val addr: %s, power: %s", val.GetOperator(), val.GetStatus().String())
	}
	t.Logf("=====================================")

	applyValidatorSetUpdates(t, ctx, appIn.StakingKeeper, -1)

	currVals = appIn.StakingKeeper.GetAllValidators(ctx)
	// print all validators addr
	for _, val := range currVals {
		t.Logf("val addr: %s, power: %s", val.GetOperator(), val.GetStatus().String())
	}
	t.Logf("=====================================")

	return addrs, valAddrs, vals
}
