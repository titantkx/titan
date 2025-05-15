package pointer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/pointer/keeper"
	"github.com/titantkx/titan/x/pointer/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the erc20Native
	for _, elem := range genState.Erc20NativeList {
		// we ignore any error record here.
		_ = k.SetErc20Native(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.Erc20NativeList = k.GetAllErc20Native(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
