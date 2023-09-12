package sims

import (
	"cosmossdk.io/math"
	sdksimtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// AddTestAddrsIncremental constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(bankKeeper bankkeeper.Keeper, ctx sdk.Context, genAddr sdk.AccAddress, accNum int, accAmt math.Int, denom string) []sdk.AccAddress {
	return addTestAddrs(bankKeeper, ctx, genAddr, accNum, accAmt, sdksimtestutil.CreateIncrementalAccounts, denom)
}

func addTestAddrs(bankKeeper bankkeeper.Keeper, ctx sdk.Context, genAddr sdk.AccAddress, accNum int, accAmt math.Int, strategy sdksimtestutil.GenerateAccountStrategy, denom string) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	initCoins := sdk.NewCoins(sdk.NewCoin(denom, accAmt))

	for _, addr := range testAddrs {
		initAccountWithCoins(bankKeeper, ctx, genAddr, addr, initCoins)
	}

	return testAddrs
}

func initAccountWithCoins(bankKeeper bankkeeper.Keeper, ctx sdk.Context, genAccounts sdk.AccAddress, addr sdk.AccAddress, coins sdk.Coins) {
	// send coin from genesis account to addr
	err := bankKeeper.SendCoins(ctx, genAccounts, addr, coins)
	if err != nil {
		panic(err)
	}
}
