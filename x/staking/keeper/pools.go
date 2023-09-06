package keeper

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// contributeBondedTokens move coins from the bonded pool module account to distribution community pool
func (k Keeper) contributeBondedTokens(ctx sdk.Context, amt math.Int) error {
	if !amt.IsPositive() {
		// skip as no coins need to be contributed
		return nil
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), amt))

	return k.distkeeper.FundCommunityPoolFromModule(ctx, coins, types.BondedPoolName)
}

// contributeNotBondedTokens move coins from the not bonded pool module account to distribution community pool
func (k Keeper) contributeNotBondedTokens(ctx sdk.Context, amt math.Int) error {
	if !amt.IsPositive() {
		// skip as no coins need to be contributed
		return nil
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), amt))

	return k.distkeeper.FundCommunityPoolFromModule(ctx, coins, types.NotBondedPoolName)
}
