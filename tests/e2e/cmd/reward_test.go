package cmd_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd/block"
	"github.com/tokenize-titan/titan/testutil/cmd/distribution"
	"github.com/tokenize-titan/titan/testutil/cmd/reward"
	"github.com/tokenize-titan/titan/testutil/cmd/staking"
	"github.com/tokenize-titan/titan/testutil/cmd/status"
	"github.com/tokenize-titan/titan/utils"
)

func MustGetRewardPoolAdmin(t testing.TB) string {
	return reward.MustGetParams(t).Authority
}

func TestSetRate(t *testing.T) {
	rewardPoolAdmin := MustGetRewardPoolAdmin(t)

	oldRate := reward.MustGetParams(t).Rate
	newRate := testutil.MakeFloat(0.1)

	reward.MustSetRate(t, rewardPoolAdmin, newRate)
	reward.MustSetRate(t, rewardPoolAdmin, oldRate)
}

func TestSetRateUnauthorized(t *testing.T) {
	t.Parallel()

	someone := MustCreateAccount(t, "1"+utils.DisplayDenom).Address

	err := reward.SetRate(someone, testutil.MakeFloat(0.1))

	require.Error(t, err)
	require.ErrorContains(t, err, "not allowed to set rate: forbidden")
}

func TestSetAuthority(t *testing.T) {
	oldAdmin := MustGetRewardPoolAdmin(t)
	newAdmin := MustCreateAccount(t, "1"+utils.DisplayDenom).Address

	reward.MustSetAuthority(t, oldAdmin, newAdmin)
	reward.MustSetAuthority(t, newAdmin, oldAdmin)
}

func TestSetAuthorityUnauthorized(t *testing.T) {
	t.Parallel()

	someone := MustCreateAccount(t, "1"+utils.DisplayDenom).Address

	err := reward.SetAuthority(someone, someone)

	require.Error(t, err)
	require.ErrorContains(t, err, "not allowed to set authority: forbidden")
}

func TestDistributeRewards(t *testing.T) {
	dornor := MustCreateAccount(t, "1001"+utils.DisplayDenom).Address

	reward.MustFundRewardPool(t, dornor, "1000"+utils.DisplayDenom)

	stakeAmount := testutil.MustParseAmount(t, "1000"+utils.DisplayDenom)

	val := MustGetValidator(t)
	bal := stakeAmount.Add(testutil.OneToken)
	del := MustCreateAccount(t, bal.GetAmount()).Address

	staking.MustDelegate(t, val, stakeAmount.String(), del)
	defer staking.MustUnbond(t, val, stakeAmount.String(), del)

	// For some reason, validator's voting power is not updated immediately after delegation
	// so we need to wait for a few blocks before rewards are correctly distributed to delegator
	startHeight := status.MustGetLatestBlockHeight(t) + 5
	endHeight := startHeight + 5

	status.MustWait(t, endHeight)

	startTime := block.MustGetBlockTime(t, startHeight)
	endTime := block.MustGetBlockTime(t, endHeight)

	rewardBefore := distribution.MustGetRewards(t, del, val, startHeight).GetBaseDenomAmount()
	rewardAfter := distribution.MustGetRewards(t, del, val, endHeight).GetBaseDenomAmount()

	interestRate := reward.MustGetParams(t).Rate
	commissionRate := staking.MustGetValidator(t, val).Commission.CommissionRates.Rate

	duration := testutil.MakeInt(endTime.Sub(startTime).Nanoseconds()).Float()
	oneYear := testutil.MakeInt((365 * 24 * time.Hour).Nanoseconds()).Float()

	interest := stakeAmount.GetBaseDenomAmount().Float().Mul(interestRate).Mul(duration).Quo(oneYear)
	commission := interest.Mul(commissionRate)

	expectedReward := interest.Sub(commission)
	actualReward := rewardAfter.Sub(rewardBefore).Float()

	if expectedReward.IsZero() { // If interest rate is zero
		require.True(t, actualReward.IsZero())
		return
	}

	diff := actualReward.Sub(expectedReward).Quo(expectedReward)
	maxDiff := testutil.MakeFloat(1e-10)

	require.Condition(
		t,
		func() bool { return diff.Abs().Cmp(maxDiff) <= 0 },
		"|%s| is not less than or equal to %s", diff, maxDiff)
}
