package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd/keys"
	"github.com/titantkx/titan/testutil/cmd/staking"
	"github.com/titantkx/titan/utils"
)

func MustGetValidator(t testing.TB) string {
	val1 := keys.MustShowAddress(t, "val1")
	return testutil.MustAccountAddressToValidatorAddress(t, val1)
}

func MustGetGlobalMinSelfDelegation(t testing.TB) testutil.Int {
	return staking.MustGetParams(t).GlobalMinSelfDelegation
}

func MustGetMinStakeAmount(_ testing.TB, minSelfDelegation testutil.Int) testutil.Int {
	stakePower := testutil.TokensToConsensusPower(minSelfDelegation)
	stakeAmount := testutil.TokensFromConsensusPower(stakePower)
	return testutil.MakeIntFromString(stakeAmount.String())
}

func MustCreateValidator(t testing.TB, stakeAmount string) (del keys.Key, val staking.Validator) {
	valPk := testutil.MustGenerateEd25519PK(t)
	del = MustCreateAccount(t, "")

	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	if stakeAmount == "" {
		stakeAmount = MustGetMinStakeAmount(t, minSelfDelegation).String() + utils.BaseDenom
	}

	bal := testutil.MustParseAmount(t, stakeAmount).Add(testutil.OneToken)
	MustAcquireMoney(t, del.Address, bal.GetAmount())

	val = staking.MustCreateValidator(t, valPk, stakeAmount, 0.1, 0.2, 0.001, minSelfDelegation, del.Address)

	return del, val
}

func MustCreateValidatorForOther(t testing.TB, stakeAmount string) (from keys.Key, del keys.Key, val staking.Validator) {
	valPk := testutil.MustGenerateEd25519PK(t)
	from = MustCreateAccount(t, "")
	del = MustCreateAccount(t, "")

	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	if stakeAmount == "" {
		stakeAmount = MustGetMinStakeAmount(t, minSelfDelegation).String() + utils.BaseDenom
	}

	bal := testutil.MustParseAmount(t, stakeAmount).Add(testutil.OneToken)
	MustAcquireMoney(t, from.Address, bal.GetAmount())

	val = staking.MustCreateValidatorForOther(t, valPk, stakeAmount, 0.1, 0.2, 0.001, minSelfDelegation, from.Address, del.Address)

	return from, del, val
}

func TestCreateValidator(t *testing.T) {
	t.Parallel()

	MustCreateValidator(t, "")
}

func TestCreateValidatorMinSelfDelegationTooLow(t *testing.T) {
	t.Parallel()

	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustCreateAccount(t, "")

	onePowerAmount := testutil.TokensFromConsensusPower(1)
	minSelfDelegation := MustGetGlobalMinSelfDelegation(t).Sub(onePowerAmount)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation).String() + utils.BaseDenom

	bal := testutil.MustParseAmount(t, stakeAmount).Add(testutil.OneToken)
	MustAcquireMoney(t, del.Address, bal.GetAmount())

	staking.MustErrCreateValidator(t, "cannot set validator min self delegation to less than global minimum", valPk, stakeAmount, 0.1, 0.2, 0.001, minSelfDelegation, del.Address)
}

func TestCreateValidatorSelfDelegationTooLow(t *testing.T) {
	t.Parallel()

	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustCreateAccount(t, "")

	onePowerAmount := testutil.TokensFromConsensusPower(1)
	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation).Sub(onePowerAmount).String() + utils.BaseDenom

	bal := testutil.MustParseAmount(t, stakeAmount).Add(testutil.OneToken)
	MustAcquireMoney(t, del.Address, bal.GetAmount())

	staking.MustErrCreateValidator(t, "validator's self delegation must be greater than their minimum self delegation", valPk, stakeAmount, 0.1, 0.2, 0.001, minSelfDelegation, del.Address)
}

func TestDelegate(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "2"+utils.DisplayDenom)
	_, val := MustCreateValidator(t, "")
	staking.MustDelegate(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)
}

func TestRedelegate(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "2"+utils.DisplayDenom)
	_, val1 := MustCreateValidator(t, "")
	_, val2 := MustCreateValidator(t, "")

	staking.MustDelegate(t, val1.OperatorAddress, "1"+utils.DisplayDenom, del.Address)
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address)
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address) // Redelegate all
}

func TestRedelegateValidatorIsJailed(t *testing.T) {
	t.Parallel()

	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation).String() + utils.BaseDenom

	del, val1 := MustCreateValidator(t, stakeAmount)
	_, val2 := MustCreateValidator(t, "")

	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "1"+utils.DisplayDenom, del.Address)

	val1 = staking.MustGetValidator(t, val1.OperatorAddress)

	require.True(t, val1.Jailed)
}

func TestRedelegateValidatorIsNotJailed(t *testing.T) {
	t.Parallel()

	redelegateAmount := testutil.MustParseAmount(t, "1"+utils.DisplayDenom)
	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation).Add(redelegateAmount.GetBaseDenomAmount()).String() + utils.BaseDenom

	del, val1 := MustCreateValidator(t, stakeAmount)
	_, val2 := MustCreateValidator(t, "")

	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, redelegateAmount.String(), del.Address)

	val1 = staking.MustGetValidator(t, val1.OperatorAddress)

	require.False(t, val1.Jailed)
}

func TestUnbond(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "2"+utils.DisplayDenom)
	_, val := MustCreateValidator(t, "")

	staking.MustDelegate(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)
	staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address)
	staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address) // Unbond all
}

func TestUnbondValidatorIsJailed(t *testing.T) {
	t.Parallel()

	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation).String() + utils.BaseDenom

	del, val := MustCreateValidator(t, stakeAmount)

	staking.MustUnbond(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)

	val = staking.MustGetValidator(t, val.OperatorAddress)

	require.True(t, val.Jailed)
}

func TestUnbondValidatorIsNotJailed(t *testing.T) {
	t.Parallel()

	unbondAmount := testutil.MustParseAmount(t, "1"+utils.DisplayDenom)
	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation).Add(unbondAmount.GetBaseDenomAmount()).String() + utils.BaseDenom

	del, val := MustCreateValidator(t, stakeAmount)

	staking.MustUnbond(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)

	val = staking.MustGetValidator(t, val.OperatorAddress)

	require.False(t, val.Jailed)
}

func TestCancelUnbound(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "2"+utils.DisplayDenom)
	_, val := MustCreateValidator(t, "")

	staking.MustDelegate(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)

	tx1 := staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address)
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.2"+utils.DisplayDenom, tx1.Height.Int64(), del.Address)

	tx2 := staking.MustUnbond(t, val.OperatorAddress, "0.7"+utils.DisplayDenom, del.Address) // Unbond all
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.7"+utils.DisplayDenom, tx2.Height.Int64(), del.Address)
}

func TestCreateValidatorForOther(t *testing.T) {
	t.Parallel()

	MustCreateValidatorForOther(t, "")
}

func TestCreateValidatorForOtherCanRedelegate(t *testing.T) {
	t.Parallel()

	_, del, val1 := MustCreateValidatorForOther(t, "")
	_, val2 := MustCreateValidator(t, "")

	MustAcquireMoney(t, del.Address, "1"+utils.DisplayDenom)
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "1"+utils.DisplayDenom, del.Address)
}

func TestCreateValidatorForOtherCanUnbond(t *testing.T) {
	t.Parallel()

	_, del, val := MustCreateValidatorForOther(t, "")

	MustAcquireMoney(t, del.Address, "1"+utils.DisplayDenom)
	staking.MustUnbond(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)
}

func TestCreateValidatorForOtherCanCancelUnbond(t *testing.T) {
	t.Parallel()

	unbondAmount := testutil.MustParseAmount(t, "1"+utils.DisplayDenom)
	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation)
	stakeAmount = stakeAmount.Add(unbondAmount.GetBaseDenomAmount()) // Make sure validator is not jailed after unbonding

	_, del, val := MustCreateValidatorForOther(t, stakeAmount.String()+utils.BaseDenom)

	MustAcquireMoney(t, del.Address, "1"+utils.DisplayDenom)
	tx := staking.MustUnbond(t, val.OperatorAddress, unbondAmount.String(), del.Address)
	staking.MustCancelUnbound(t, val.OperatorAddress, unbondAmount.String(), tx.Height.Int64(), del.Address)
}

func TestDelegateForOther(t *testing.T) {
	t.Parallel()

	del1 := MustCreateAccount(t, "2"+utils.DisplayDenom)
	del2 := MustCreateAccount(t, "")
	_, val := MustCreateValidator(t, "")

	staking.MustDelegateForOther(t, val.OperatorAddress, "1"+utils.DisplayDenom, del1.Address, del2.Address)
}

func TestDelegateForOtherCanRedelegate(t *testing.T) {
	t.Parallel()

	del1 := MustCreateAccount(t, "2"+utils.DisplayDenom)
	del2 := MustCreateAccount(t, "1"+utils.DisplayDenom)
	_, val1 := MustCreateValidator(t, "")
	_, val2 := MustCreateValidator(t, "")

	staking.MustDelegateForOther(t, val1.OperatorAddress, "1"+utils.DisplayDenom, del1.Address, del2.Address)

	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address)
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address) // Redelegate all
}

func TestDelegateForOtherCanUnbond(t *testing.T) {
	t.Parallel()

	del1 := MustCreateAccount(t, "2"+utils.DisplayDenom)
	del2 := MustCreateAccount(t, "1"+utils.DisplayDenom)
	_, val := MustCreateValidator(t, "")

	staking.MustDelegateForOther(t, val.OperatorAddress, "1"+utils.DisplayDenom, del1.Address, del2.Address)

	staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address)
	staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address) // Unbond all
}

func TestDelegateForOtherCanCancelUnbond(t *testing.T) {
	t.Parallel()

	del1 := MustCreateAccount(t, "2"+utils.DisplayDenom)
	del2 := MustCreateAccount(t, "1"+utils.DisplayDenom)
	_, val := MustCreateValidator(t, "")

	staking.MustDelegateForOther(t, val.OperatorAddress, "1"+utils.DisplayDenom, del1.Address, del2.Address)

	tx1 := staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address)
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, tx1.Height.Int64(), del2.Address)
	tx2 := staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address) // Unbond all
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, tx2.Height.Int64(), del2.Address)
}
