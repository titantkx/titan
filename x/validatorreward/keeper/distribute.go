package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/validatorreward/types"
)

const timeYear = time.Hour * 24 * 365

func (k Keeper) DistributeTokens(ctx sdk.Context, totalPreviousPower int64) {
	// get current balance of reward pool
	validatorRewardAccount := k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
	currentBalance := k.bankKeeper.GetBalance(ctx, validatorRewardAccount.GetAddress(), utils.BondDenom)

	// if current balance is zero, ignore distribution
	if currentBalance.IsZero() {
		k.SetLastDistributeTime(ctx, ctx.BlockHeader().Time)
		return
	}

	lastDistributeTime := k.GetLastDistributeTime(ctx)
	// if lastDistributeTime is zero, must wait for next block
	if lastDistributeTime.IsZero() {
		k.SetLastDistributeTime(ctx, ctx.BlockHeader().Time)
		return
	}

	apyRate := k.GetParams(ctx).Rate

	// Calculate Duration since last distribution
	duration := ctx.BlockHeader().Time.Sub(lastDistributeTime)
	yearDurationInNanoseconds := int64(timeYear)
	durationInYearFraction := sdk.NewDec(duration.Nanoseconds()).Quo(sdk.NewDec(yearDurationInNanoseconds))

	totalPreviousPowerInDecCoin := sdk.NewDecCoin(utils.BondDenom, sdk.NewInt(totalPreviousPower).Mul(sdk.DefaultPowerReduction))
	totalPreviousPowerInDecCoins := sdk.NewDecCoins(totalPreviousPowerInDecCoin)

	// Calculate reward amount
	rewardPerYearAmountDecCoin := totalPreviousPowerInDecCoins.MulDecTruncate(apyRate)
	rewardAmountDecCoin := rewardPerYearAmountDecCoin.MulDecTruncate(durationInYearFraction)
	rewardAmount, _ := rewardAmountDecCoin.TruncateDecimal()

	rewardAmount = sdk.NewCoins(currentBalance).Min(rewardAmount)

	k.SetLastDistributeTime(ctx, ctx.BlockHeader().Time)

	k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.ValidatorRewardCollectorName, rewardAmount)
}
