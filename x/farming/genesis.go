package farming

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/farming/keeper"
	"github.com/titantkx/titan/x/farming/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the farm
	for _, elem := range genState.FarmList {
		k.SetFarm(ctx, elem)
	}
	// Set all the stakingInfo
	for _, elem := range genState.StakingInfoList {
		k.SetStakingInfo(ctx, elem)
	}
	// Set if defined
	if genState.DistributionInfo != nil {
		k.SetDistributionInfo(ctx, *genState.DistributionInfo)
	}
	// Set all the reward
	for _, elem := range genState.RewardList {
		k.SetReward(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.FarmList = k.GetAllFarm(ctx)
	genesis.StakingInfoList = k.GetAllStakingInfo(ctx, "")
	// Get distributionInfo
	distributionInfo, found := k.GetDistributionInfo(ctx)
	if found {
		genesis.DistributionInfo = &distributionInfo
	}
	genesis.RewardList = k.GetAllReward(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
