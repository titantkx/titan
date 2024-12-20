package farming

import (
	"time"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/farming/keeper"
	"github.com/titantkx/titan/x/farming/types"
)

// BeginBlocker add farming reward to stakers
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	var lastDistributionTime time.Time

	distributionInfo, found := k.GetDistributionInfo(ctx)
	if found {
		lastDistributionTime = distributionInfo.LastDistributionTime
	}

	farms := k.GetAllFarm(ctx)

	for _, farm := range farms {
		stakingInfos := k.GetAllStakingInfo(ctx, farm.Token)
		totalStakedAmount := math.NewInt(0)

		for _, stakingInfo := range stakingInfos {
			totalStakedAmount = totalStakedAmount.Add(stakingInfo.Amount)
		}

		pendingRewards := make([]*types.FarmReward, 0, len(farm.Rewards))

		for _, reward := range farm.Rewards {
			if !req.Header.Time.After(reward.StartTime) {
				pendingRewards = append(pendingRewards, reward)
				continue
			}

			var totalDistributedRewardAmount sdk.Coins

			if !reward.EndTime.After(req.Header.Time) {
				totalDistributedRewardAmount = reward.Amount
			} else {
				var beginRewardDistributionTime time.Time

				if lastDistributionTime.After(reward.StartTime) {
					beginRewardDistributionTime = lastDistributionTime
				} else {
					beginRewardDistributionTime = reward.StartTime
				}

				totalDistributedRewardAmount = reward.Amount.
					MulInt(sdk.NewInt(int64(req.Header.Time.Sub(beginRewardDistributionTime)))).
					QuoInt(sdk.NewInt(int64(reward.EndTime.Sub(beginRewardDistributionTime))))
			}

			for _, stakingInfo := range stakingInfos {
				distributedRewardAmount := totalDistributedRewardAmount.
					MulInt(stakingInfo.Amount).
					QuoInt(totalStakedAmount)

				reward.Amount = reward.Amount.Sub(distributedRewardAmount...)
				k.AddReward(ctx, stakingInfo.Staker, distributedRewardAmount)
			}

			if reward.EndTime.After(req.Header.Time) {
				pendingRewards = append(pendingRewards, reward)
			} else if !reward.Amount.IsZero() {
				// Refund unused rewards to the sender
				k.AddReward(ctx, reward.Sender, reward.Amount)
			}
		}

		farm.Rewards = pendingRewards

		if len(farm.Rewards) == 0 {
			k.RemoveFarm(ctx, farm.Token)
		} else {
			k.SetFarm(ctx, farm)
		}
	}

	distributionInfo.LastDistributionTime = req.Header.Time
	k.SetDistributionInfo(ctx, distributionInfo)
}
