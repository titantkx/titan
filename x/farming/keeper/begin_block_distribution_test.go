package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/farming/keeper"
	"github.com/titantkx/titan/x/farming/types"
)

func checkRewardAmount(t testing.TB, k *keeper.Keeper, ctx sdk.Context, address string, amount sdk.Coins) {
	reward, found := k.GetReward(ctx, address)

	if amount.IsZero() {
		require.False(t, found)
	} else {
		require.True(t, found)
		require.True(t, amount.IsEqual(reward.Amount))
	}
}

func checkPendingRewards(t testing.TB, k *keeper.Keeper, ctx sdk.Context, token string, amount sdk.Coins) {
	farm, found := k.GetFarm(ctx, token)

	if amount.IsZero() {
		require.False(t, found)
	} else {
		require.True(t, found)
		require.NotEmpty(t, farm.Rewards)
		require.True(t, amount.IsEqual(farm.Rewards[0].Amount))
	}
}

func TestDistributeRewards(t *testing.T) {
	k, ctx := keepertest.FarmingKeeper(t)

	startTime := time.Now()
	endTime := startTime.Add(10 * time.Hour)

	sender := sample.AccAddress().String()
	farmer1 := sample.AccAddress().String()
	farmer2 := sample.AccAddress().String()

	k.SetFarm(ctx, types.Farm{
		Token: "btc",
		Rewards: []*types.FarmReward{
			{
				Sender:    sender,
				Amount:    utils.NewCoins("100tkx"),
				EndTime:   endTime,
				StartTime: startTime,
			},
		},
	})
	k.SetStakingInfo(ctx, types.StakingInfo{
		Token:  "btc",
		Staker: farmer1,
		Amount: math.NewInt(2),
	})
	k.SetStakingInfo(ctx, types.StakingInfo{
		Token:  "btc",
		Staker: farmer2,
		Amount: math.NewInt(3),
	})

	k.DistributeRewards(ctx, startTime)

	checkRewardAmount(t, k, ctx, farmer1, sdk.NewCoins())
	checkRewardAmount(t, k, ctx, farmer2, sdk.NewCoins())
	checkPendingRewards(t, k, ctx, "btc", utils.NewCoins("100tkx"))

	k.DistributeRewards(ctx, startTime.Add(1*time.Hour))

	checkRewardAmount(t, k, ctx, farmer1, utils.NewCoins("4tkx"))
	checkRewardAmount(t, k, ctx, farmer2, utils.NewCoins("6tkx"))
	checkPendingRewards(t, k, ctx, "btc", utils.NewCoins("90tkx"))

	k.RemoveStakingInfo(ctx, "btc", farmer1)
	k.RemoveStakingInfo(ctx, "btc", farmer2)

	k.DistributeRewards(ctx, startTime.Add(2*time.Hour))

	checkRewardAmount(t, k, ctx, farmer1, utils.NewCoins("4tkx"))
	checkRewardAmount(t, k, ctx, farmer2, utils.NewCoins("6tkx"))
	checkPendingRewards(t, k, ctx, "btc", utils.NewCoins("90tkx"))

	k.DistributeRewards(ctx, endTime)

	checkRewardAmount(t, k, ctx, farmer1, utils.NewCoins("4tkx"))
	checkRewardAmount(t, k, ctx, farmer2, utils.NewCoins("6tkx"))
	checkRewardAmount(t, k, ctx, sender, utils.NewCoins("90tkx"))
	checkPendingRewards(t, k, ctx, "btc", sdk.NewCoins())
}
