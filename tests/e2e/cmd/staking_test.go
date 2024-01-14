package cmd_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd/staking"
	"github.com/tokenize-titan/titan/utils"
)

func MustCreateValidator(t testing.TB) staking.Validator {
	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustAddKey(t)
	MustAcquireMoney(t, del.Address, "1000"+utils.DisplayDenom)
	stakingParams := staking.MustGetParams(t)
	stakePower := sdk.TokensToConsensusPower(stakingParams.GlobalMinSelfDelegation, sdk.DefaultPowerReduction) + 1
	stakeAmount := sdk.TokensFromConsensusPower(stakePower, sdk.DefaultPowerReduction)
	return staking.MustCreateValidator(t, valPk, fmt.Sprintf("%s%s", stakeAmount.String(), utils.BaseDenom), 0.1, 0.2, 0.001, stakingParams.GlobalMinSelfDelegation, del.Address)
}

func MustErrCreateValidator(t testing.TB) {
	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustAddKey(t)
	MustAcquireMoney(t, del.Address, "1000"+utils.DisplayDenom)
	stakingParams := staking.MustGetParams(t)
	if !(stakingParams.GlobalMinSelfDelegation.Int64() > 1) {
		// This test requires that the global min self delegation is greater than 1
		return
	}

	onePowerAmount := sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)
	twoPowerAmount := sdk.TokensFromConsensusPower(2, sdk.DefaultPowerReduction)

	stakeAmount := stakingParams.GlobalMinSelfDelegation.Add(onePowerAmount)

	staking.MustErrCreateValidator(t, "cannot set validator min self delegation to less than global minimum", valPk, fmt.Sprintf("%s%s", stakeAmount.String(), utils.BaseDenom), 0.1, 0.2, 0.001, stakingParams.GlobalMinSelfDelegation.Sub(onePowerAmount), del.Address)

	staking.MustErrCreateValidator(t, "", valPk, fmt.Sprintf("%s%s", stakeAmount.String(), utils.BaseDenom), 0.1, 0.2, 0.001, stakingParams.GlobalMinSelfDelegation.Add(twoPowerAmount), del.Address)
}

func TestCreateValidator(t *testing.T) {
	t.Parallel()

	MustCreateValidator(t)
	MustErrCreateValidator(t)
}

func TestDelegate(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "1000"+utils.DisplayDenom)
	val := MustCreateValidator(t)
	staking.MustDelegate(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)
}

func TestRedelegate(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "1000"+utils.DisplayDenom)
	val1 := MustCreateValidator(t)
	val2 := MustCreateValidator(t)

	staking.MustDelegate(t, val1.OperatorAddress, "1"+utils.DisplayDenom, del.Address)
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address)
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address) // Redelegate all
}

func TestUnbond(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "1000"+utils.DisplayDenom)
	val := MustCreateValidator(t)

	staking.MustDelegate(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)
	staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address)
	staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address) // Unbond all
}

func TestCancelUnbound(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "1000"+utils.DisplayDenom)
	val := MustCreateValidator(t)

	staking.MustDelegate(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)

	tx1 := staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address)
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.2"+utils.DisplayDenom, tx1.Height.Int64(), del.Address)

	tx2 := staking.MustUnbond(t, val.OperatorAddress, "0.7"+utils.DisplayDenom, del.Address) // Unbond all
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.7"+utils.DisplayDenom, tx2.Height.Int64(), del.Address)
}
