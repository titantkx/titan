package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd/keys"
	"github.com/tokenize-titan/titan/testutil/cmd/staking"
	"github.com/tokenize-titan/titan/utils"
)

func MustGetGlobalMinSelfDelegation(t testing.TB) testutil.BigInt {
	return staking.MustGetParams(t).GlobalMinSelfDelegation
}

func MustGetMinStakeAmount(t testing.TB, minSelfDelegation testutil.BigInt) testutil.BigInt {
	stakePower := testutil.TokensToConsensusPower(minSelfDelegation)
	stakeAmount := testutil.TokensFromConsensusPower(stakePower)
	return testutil.MakeBigIntFromString(stakeAmount.String())
}

func MustCreateValidator(t testing.TB, stakeAmount string) staking.Validator {
	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustCreateAccount(t, "10000"+utils.DisplayDenom)

	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	if stakeAmount == "" {
		stakeAmount = MustGetMinStakeAmount(t, minSelfDelegation).String() + utils.BaseDenom
	}

	return staking.MustCreateValidator(t, valPk, stakeAmount, 0.1, 0.2, 0.001, minSelfDelegation, del.Address)
}

func MustCreateValidatorForOther(t testing.TB, stakeAmount string) (from keys.Key, del keys.Key, val staking.Validator) {
	valPk := testutil.MustGenerateEd25519PK(t)
	from = MustCreateAccount(t, "10000"+utils.DisplayDenom)
	del = MustCreateAccount(t, "")

	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	if stakeAmount == "" {
		stakeAmount = MustGetMinStakeAmount(t, minSelfDelegation).String() + utils.BaseDenom
	}

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
	del := MustCreateAccount(t, "10000"+utils.DisplayDenom)

	onePowerAmount := testutil.TokensFromConsensusPower(1)
	minSelfDelegation := MustGetGlobalMinSelfDelegation(t).Sub(onePowerAmount)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation)

	staking.MustErrCreateValidator(t, "cannot set validator min self delegation to less than global minimum", valPk, stakeAmount.String()+utils.BaseDenom, 0.1, 0.2, 0.001, minSelfDelegation, del.Address)
}

func TestCreateValidatorSelfDelegationTooLow(t *testing.T) {
	t.Parallel()

	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustCreateAccount(t, "10000"+utils.DisplayDenom)

	onePowerAmount := testutil.TokensFromConsensusPower(1)
	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation).Sub(onePowerAmount)

	staking.MustErrCreateValidator(t, "validator's self delegation must be greater than their minimum self delegation", valPk, stakeAmount.String()+utils.BaseDenom, 0.1, 0.2, 0.001, minSelfDelegation, del.Address)
}

func TestDelegate(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "2"+utils.DisplayDenom)
	val := MustCreateValidator(t, "")
	staking.MustDelegate(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)
}

func TestRedelegate(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "2"+utils.DisplayDenom)
	val1 := MustCreateValidator(t, "")
	val2 := MustCreateValidator(t, "")

	staking.MustDelegate(t, val1.OperatorAddress, "1"+utils.DisplayDenom, del.Address)
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address)
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address) // Redelegate all
}

func TestRedelegateValidatorIsJailed(t *testing.T) {
	t.Parallel()

	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustCreateAccount(t, "10000"+utils.DisplayDenom)

	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation).String() + utils.BaseDenom

	val1 := staking.MustCreateValidator(t, valPk, stakeAmount, 0.1, 0.2, 0.001, minSelfDelegation, del.Address)
	val2 := MustCreateValidator(t, "")

	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "1"+utils.DisplayDenom, del.Address)

	val1 = staking.MustGetValidator(t, val1.OperatorAddress)

	require.True(t, val1.Jailed)
}

func TestRedelegateValidatorIsNotJailed(t *testing.T) {
	t.Parallel()

	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustCreateAccount(t, "10000"+utils.DisplayDenom)

	redelegateAmount := testutil.MustParseAmount(t, "1"+utils.DisplayDenom)
	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation).Add(redelegateAmount.GetBaseDenomAmount())

	val1 := staking.MustCreateValidator(t, valPk, stakeAmount.String()+utils.BaseDenom, 0.1, 0.2, 0.001, minSelfDelegation, del.Address)
	val2 := MustCreateValidator(t, "")

	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, redelegateAmount.String(), del.Address)

	val1 = staking.MustGetValidator(t, val1.OperatorAddress)

	require.False(t, val1.Jailed)
}

func TestUnbond(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "2"+utils.DisplayDenom)
	val := MustCreateValidator(t, "")

	staking.MustDelegate(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)
	staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address)
	staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del.Address) // Unbond all
}

func TestUnbondValidatorIsJailed(t *testing.T) {
	t.Parallel()

	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustCreateAccount(t, "10000"+utils.DisplayDenom)

	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation).String() + utils.BaseDenom

	val := staking.MustCreateValidator(t, valPk, stakeAmount, 0.1, 0.2, 0.001, minSelfDelegation, del.Address)

	staking.MustUnbond(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)

	val = staking.MustGetValidator(t, val.OperatorAddress)

	require.True(t, val.Jailed)
}

func TestUnbondValidatorIsNotJailed(t *testing.T) {
	t.Parallel()

	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustCreateAccount(t, "10000"+utils.DisplayDenom)

	unbondAmount := testutil.MustParseAmount(t, "1"+utils.DisplayDenom)
	minSelfDelegation := MustGetGlobalMinSelfDelegation(t)
	stakeAmount := MustGetMinStakeAmount(t, minSelfDelegation).Add(unbondAmount.GetBaseDenomAmount())

	val := staking.MustCreateValidator(t, valPk, stakeAmount.String()+utils.BaseDenom, 0.1, 0.2, 0.001, minSelfDelegation, del.Address)

	staking.MustUnbond(t, val.OperatorAddress, "1"+utils.DisplayDenom, del.Address)

	val = staking.MustGetValidator(t, val.OperatorAddress)

	require.False(t, val.Jailed)
}

func TestCancelUnbound(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "2"+utils.DisplayDenom)
	val := MustCreateValidator(t, "")

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
	val2 := MustCreateValidator(t, "")

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
	val := MustCreateValidator(t, "")

	staking.MustDelegateForOther(t, val.OperatorAddress, "1"+utils.DisplayDenom, del1.Address, del2.Address)
}

func TestDelegateForOtherCanRedelegate(t *testing.T) {
	t.Parallel()

	del1 := MustCreateAccount(t, "2"+utils.DisplayDenom)
	del2 := MustCreateAccount(t, "1"+utils.DisplayDenom)
	val1 := MustCreateValidator(t, "")
	val2 := MustCreateValidator(t, "")

	staking.MustDelegateForOther(t, val1.OperatorAddress, "1"+utils.DisplayDenom, del1.Address, del2.Address)

	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address)
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address) // Redelegate all
}

func TestDelegateForOtherCanUnbond(t *testing.T) {
	t.Parallel()

	del1 := MustCreateAccount(t, "2"+utils.DisplayDenom)
	del2 := MustCreateAccount(t, "1"+utils.DisplayDenom)
	val := MustCreateValidator(t, "")

	staking.MustDelegateForOther(t, val.OperatorAddress, "1"+utils.DisplayDenom, del1.Address, del2.Address)

	staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address)
	staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address) // Unbond all
}

func TestDelegateForOtherCanCancelUnbond(t *testing.T) {
	t.Parallel()

	del1 := MustCreateAccount(t, "2"+utils.DisplayDenom)
	del2 := MustCreateAccount(t, "1"+utils.DisplayDenom)
	val := MustCreateValidator(t, "")

	staking.MustDelegateForOther(t, val.OperatorAddress, "1"+utils.DisplayDenom, del1.Address, del2.Address)

	tx1 := staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address)
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, tx1.Height.Int64(), del2.Address)
	tx2 := staking.MustUnbond(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, del2.Address) // Unbond all
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.5"+utils.DisplayDenom, tx2.Height.Int64(), del2.Address)
}
