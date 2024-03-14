package nftmint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/nftmint/keeper"
	"github.com/titantkx/titan/x/nftmint/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set systemInfo
	k.SetSystemInfo(ctx, genState.SystemInfo)
	// Set all the mintingInfo
	for _, elem := range genState.MintingInfoList {
		k.SetMintingInfo(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	err := k.SetParams(ctx, genState.Params)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Get all systemInfo
	systemInfo, found := k.GetSystemInfo(ctx)
	if found {
		genesis.SystemInfo = systemInfo
	}
	genesis.MintingInfoList = k.GetAllMintingInfo(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
