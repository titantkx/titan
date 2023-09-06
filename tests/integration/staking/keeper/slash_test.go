package keeper_test

import (
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkstakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/tokenize-titan/titan/app"
)

// bootstrapSlashTest creates 3 validators and bootstrap the app.
func (s *IntegrationTestSuite) bootstrapSlashTest(power int64) (*app.App, sdk.Context, []sdk.AccAddress, []sdk.ValAddress) {
	require := s.Require()

	addrDels, addrVals := generateAddresses(s.app, s.ctx, s.genAddr, 100)

	amt := s.app.StakingKeeper.TokensFromConsensusPower(s.ctx, power)
	totalSupply := sdk.NewCoins(sdk.NewCoin(s.app.StakingKeeper.BondDenom(s.ctx), amt.MulRaw(int64(len(addrDels)))))

	notBondedPool := s.app.StakingKeeper.GetNotBondedPool(s.ctx)
	require.NoError(s.app.BankKeeper.SendCoinsFromAccountToModule(s.ctx, s.genAddr, notBondedPool.GetName(), totalSupply))

	s.app.AccountKeeper.SetModuleAccount(s.ctx, notBondedPool)

	numVals := int64(3)
	bondedCoins := sdk.NewCoins(sdk.NewCoin(s.app.StakingKeeper.BondDenom(s.ctx), amt.MulRaw(numVals)))
	bondedPool := s.app.StakingKeeper.GetBondedPool(s.ctx)

	// set bonded pool balance
	s.app.AccountKeeper.SetModuleAccount(s.ctx, bondedPool)
	require.NoError(s.app.BankKeeper.SendCoinsFromAccountToModule(s.ctx, s.genAddr, bondedPool.GetName(), bondedCoins))

	for i := int64(0); i < numVals; i++ {
		validator := testutil.NewValidator(s.T(), addrVals[i], PKs[i])
		validator, _ = validator.AddTokensFromDel(amt)
		validator = sdkstakingkeeper.TestingUpdateValidator(s.app.StakingKeeper.Keeper, s.ctx, validator, true)
		s.app.StakingKeeper.SetValidatorByConsAddr(s.ctx, validator)
	}

	return s.app, s.ctx, addrDels, addrVals
}

// tests slashUnbondingDelegation
func (s *IntegrationTestSuite) TestSlashUnbondingDelegation() {
	require := s.Require()

	app, ctx, addrDels, addrVals := s.bootstrapSlashTest(10)

	fraction := sdk.NewDecWithPrec(5, 1)

	// set an unbonding delegation with expiration timestamp (beyond which the
	// unbonding delegation shouldn't be slashed)
	ubd := types.NewUnbondingDelegation(addrDels[0], addrVals[0], 0,
		time.Unix(5, 0), sdk.NewInt(10), 0)

	app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)

	// unbonding started prior to the infraction height, stakw didn't contribute
	slashAmount := app.StakingKeeper.SlashUnbondingDelegation(ctx, ubd, 1, fraction)
	require.True(slashAmount.Equal(sdk.NewInt(0)))

	// after the expiration time, no longer eligible for slashing
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: time.Unix(10, 0)})
	app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
	slashAmount = app.StakingKeeper.SlashUnbondingDelegation(ctx, ubd, 0, fraction)
	require.True(slashAmount.Equal(sdk.NewInt(0)))

	// check community pool before slash
	oldCommunityPoolBalance := app.DistrKeeper.GetFeePoolCommunityCoins(ctx).AmountOf(app.StakingKeeper.BondDenom(ctx))
	require.Equal(sdk.NewDec(0), oldCommunityPoolBalance)

	// test valid slash, before expiration timestamp and to which stake contributed
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	oldUnbondedPoolBalances := app.BankKeeper.GetAllBalances(ctx, notBondedPool.GetAddress())
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: time.Unix(0, 0)})
	app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
	slashAmount = app.StakingKeeper.SlashUnbondingDelegation(ctx, ubd, 0, fraction)
	require.True(slashAmount.Equal(sdk.NewInt(5)))
	ubd, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.True(found)
	require.Len(ubd.Entries, 1)

	// initial balance unchanged
	require.Equal(sdk.NewInt(10), ubd.Entries[0].InitialBalance)

	// balance decreased
	require.Equal(sdk.NewInt(5), ubd.Entries[0].Balance)
	newUnbondedPoolBalances := app.BankKeeper.GetAllBalances(ctx, notBondedPool.GetAddress())
	diffTokens := oldUnbondedPoolBalances.Sub(newUnbondedPoolBalances...)
	require.True(diffTokens.AmountOf(app.StakingKeeper.BondDenom(ctx)).Equal(sdk.NewInt(5)))

	// check community pool
	newCommunityPoolBalance := app.DistrKeeper.GetFeePoolCommunityCoins(ctx).AmountOf(app.StakingKeeper.BondDenom(ctx))
	require.Equal(sdk.NewDec(5), newCommunityPoolBalance)
}
