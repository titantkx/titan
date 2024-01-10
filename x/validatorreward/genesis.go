package validatorreward

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenize-titan/titan/x/validatorreward/keeper"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if genState.LastDistributeTime != nil {
		k.SetLastDistributeTime(ctx, *genState.LastDistributeTime)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	lastDistributeTime := k.GetLastDistributeTime(ctx)
	if !lastDistributeTime.IsZero() {
		genesis.LastDistributeTime = &lastDistributeTime
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
