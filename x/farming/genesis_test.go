package farming_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/farming"
	"github.com/titantkx/titan/x/farming/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		FarmList: []types.Farm{
			{
				Token: "0",
				Rewards: []*types.FarmReward{
					{
						Sender:    sample.AccAddress().String(),
						Amount:    utils.NewCoins("1000btc"),
						EndTime:   time.Now().Add(1 * time.Hour),
						StartTime: time.Now(),
					},
				},
			},
			{
				Token: "1",
			},
		},
		StakingInfoList: []types.StakingInfo{
			{
				Token:  "0",
				Staker: "0",
				Amount: math.NewInt(100),
			},
			{
				Token:  "1",
				Staker: "1",
				Amount: math.NewInt(1000),
			},
		},
		DistributionInfo: &types.DistributionInfo{
			LastDistributionTime: time.Now(),
		},
		RewardList: []types.Reward{
			{
				Farmer: "0",
			},
			{
				Farmer: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.FarmingKeeper(t)
	farming.InitGenesis(ctx, *k, genesisState)
	got := farming.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.FarmList, got.FarmList)
	require.ElementsMatch(t, genesisState.StakingInfoList, got.StakingInfoList)
	require.True(t, genesisState.DistributionInfo.LastDistributionTime.Equal(got.DistributionInfo.LastDistributionTime))
	require.ElementsMatch(t, genesisState.RewardList, got.RewardList)
	// this line is used by starport scaffolding # genesis/test/assert
}
