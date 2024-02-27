package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd/bank"
	"github.com/tokenize-titan/titan/testutil/cmd/distribution"
	"github.com/tokenize-titan/titan/testutil/cmd/keys"
	"github.com/tokenize-titan/titan/testutil/cmd/reward"
	"github.com/tokenize-titan/titan/testutil/cmd/staking"
	"github.com/tokenize-titan/titan/testutil/cmd/status"
	"github.com/tokenize-titan/titan/utils"
)

func TestDistributeTransactionFees(t *testing.T) {
	rewardPoolAdmin := MustGetRewardPoolAdmin(t)
	oldInterestRate := reward.MustGetParams(t).Rate
	reward.MustSetRate(t, rewardPoolAdmin, testutil.MakeFloat(0)) // Set interest rate to zero
	defer reward.MustSetRate(t, rewardPoolAdmin, oldInterestRate)

	// Wait for one block to distribute previous transaction's fee
	startHeight := status.MustGetLatestBlockHeight(t) + 1
	status.MustWait(t, startHeight)

	// Make some transactions and collect their fees
	txFees := testutil.MakeInt(0)
	for i := 0; i < 5; i++ {
		faucet := keys.MustShowAddress(t, "faucet")
		receiver := MustAddKey(t).Address
		txr := bank.MustSend(t, faucet, receiver, "1"+utils.DisplayDenom)
		txFees = txFees.Add(txr.MustGetDeductFeeAmount(t))
	}

	// Wait for one block to distribute transaction's fee
	endHeight := status.MustGetLatestBlockHeight(t) + 1
	status.MustWait(t, endHeight)

	del := keys.MustShowAddress(t, "val1")
	val := testutil.MustAccountAddressToValidatorAddress(t, del)

	valShares := staking.MustGetValidator(t, val).Tokens.Float()
	totalShares := staking.MustGetStakingPool(t).BondedTokens.Float()

	rewardBefore := distribution.MustGetRewards(t, del, val, startHeight).GetBaseDenomAmount()
	rewardAfter := distribution.MustGetRewards(t, del, val, endHeight).GetBaseDenomAmount()

	communityTax := distribution.MustGetParams(t).CommunityTax
	commissionRate := staking.MustGetValidator(t, val).Commission.CommissionRates.Rate

	expectedReward := txFees.Float().Mul(valShares).Quo(totalShares).Mul(testutil.MakeFloat(1).Sub(communityTax)).Mul(testutil.MakeFloat(1).Sub(commissionRate))
	actualReward := rewardAfter.Sub(rewardBefore).Float()

	require.Equal(t, expectedReward.String(), actualReward.String())
}
