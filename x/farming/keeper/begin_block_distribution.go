package keeper

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/farming/types"
)

func (k Keeper) DistributeRewards(ctx sdk.Context, blockTime time.Time) {
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
			if !blockTime.After(reward.StartTime) {
				pendingRewards = append(pendingRewards, reward)
				continue
			}

			var totalDistributedRewardAmount sdk.Coins

			if !reward.EndTime.After(blockTime) {
				totalDistributedRewardAmount = reward.Amount
			} else {
				var beginDistributionTime time.Time

				if lastDistributionTime.After(reward.StartTime) {
					beginDistributionTime = lastDistributionTime
				} else {
					beginDistributionTime = reward.StartTime
				}

				totalDistributedRewardAmount = reward.Amount.
					MulInt(sdk.NewInt(int64(blockTime.Sub(beginDistributionTime)))).
					QuoInt(sdk.NewInt(int64(reward.EndTime.Sub(beginDistributionTime))))
			}

			for _, stakingInfo := range stakingInfos {
				distributedRewardAmount := totalDistributedRewardAmount.
					MulInt(stakingInfo.Amount).
					QuoInt(totalStakedAmount)

				reward.Amount = reward.Amount.Sub(distributedRewardAmount...)
				k.AddReward(ctx, stakingInfo.Staker, distributedRewardAmount)
			}

			if reward.EndTime.After(blockTime) {
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

	distributionInfo.LastDistributionTime = blockTime
	k.SetDistributionInfo(ctx, distributionInfo)
}
