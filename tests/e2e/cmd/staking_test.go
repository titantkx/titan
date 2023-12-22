package cmd_test

import (
	"testing"

	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd/staking"
	"github.com/tokenize-titan/titan/utils"
)

func MustCreateValidator(t testing.TB) staking.Validator {
	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustAddKey(t)
	MustAcquireMoney(t, del.Address, "1000"+utils.DisplayDenom)
	return staking.MustCreateValidator(t, valPk, "1"+utils.DisplayDenom, 0.1, 0.2, 0.001, 1, del.Address)
}

func TestCreateValidator(t *testing.T) {
	t.Parallel()

	MustCreateValidator(t)
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
