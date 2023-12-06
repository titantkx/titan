package cmd_test

import (
	"testing"

	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd/staking"
)

func MustCreateValidator(t testing.TB) staking.Validator {
	valPk := testutil.MustGenerateEd25519PK(t)
	del := MustAddKey(t)
	MustAcquireMoney(t, del.Address, "1tkx")
	return staking.MustCreateValidator(t, valPk, "0.001tkx", 0.1, 0.2, 0.001, 1, del.Address)
}

func TestCreateValidator(t *testing.T) {
	t.Parallel()

	MustCreateValidator(t)
}

func TestDelegate(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "1tkx")
	val := MustCreateValidator(t)
	staking.MustDelegate(t, val.OperatorAddress, "0.001tkx", del.Address)
}

func TestRedelegate(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "1tkx")
	val1 := MustCreateValidator(t)
	val2 := MustCreateValidator(t)

	staking.MustDelegate(t, val1.OperatorAddress, "0.001tkx", del.Address)
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.0005tkx", del.Address)
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.0005tkx", del.Address) // Redelegate all
}

func TestUnbond(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "1tkx")
	val := MustCreateValidator(t)

	staking.MustDelegate(t, val.OperatorAddress, "0.001tkx", del.Address)
	staking.MustUnbond(t, val.OperatorAddress, "0.0005tkx", del.Address)
	staking.MustUnbond(t, val.OperatorAddress, "0.0005tkx", del.Address) // Unbond all
}

func TestCancelUnbound(t *testing.T) {
	t.Parallel()

	del := MustCreateAccount(t, "1tkx")
	val := MustCreateValidator(t)

	staking.MustDelegate(t, val.OperatorAddress, "0.001tkx", del.Address)

	tx1 := staking.MustUnbond(t, val.OperatorAddress, "0.0005tkx", del.Address)
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.0002tkx", tx1.Height.Int64(), del.Address)

	tx2 := staking.MustUnbond(t, val.OperatorAddress, "0.0007tkx", del.Address) // Unbond all
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.0007tkx", tx2.Height.Int64(), del.Address)
}
