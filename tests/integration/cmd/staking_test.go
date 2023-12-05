package cmd_test

import (
	"testing"

	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd/keys"
	"github.com/tokenize-titan/titan/testutil/cmd/staking"
)

func createValidator(t testing.TB) staking.Validator {
	// Generate a fake public key for the validator node
	valPk := testutil.MustGenerateEd25519PK(t)
	// Create delegator account
	delName := testutil.GetName()
	defer testutil.PutName(delName)
	defer keys.MustDelete(t, delName)
	del := keys.MustAdd(t, delName)
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, del.Address, "1tkx")
	// Stake 0.001tkx
	return staking.MustCreateValidator(t, valPk, "0.001tkx", 0.1, 0.2, 0.001, 1, del.Address)
}

func TestCreateValidator(t *testing.T) {
	t.Parallel()

	createValidator(t)
}

func TestDelegate(t *testing.T) {
	t.Parallel()

	// Create validator
	val := createValidator(t)
	// Create delegator account
	delName := testutil.GetName()
	defer testutil.PutName(delName)
	defer keys.MustDelete(t, delName)
	del := keys.MustAdd(t, delName)
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, del.Address, "1tkx")
	// Delegate 0.001tkx
	staking.MustDelegate(t, val.OperatorAddress, "0.001tkx", del.Address)
}

func TestRedelegate(t *testing.T) {
	t.Parallel()

	// Create validator 1
	val1 := createValidator(t)
	// Create validator 2
	val2 := createValidator(t)
	// Create delegator account
	delName := testutil.GetName()
	defer testutil.PutName(delName)
	defer keys.MustDelete(t, delName)
	del := keys.MustAdd(t, delName)
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, del.Address, "1tkx")
	// Delegate 0.001tkx to validator 1
	staking.MustDelegate(t, val1.OperatorAddress, "0.001tkx", del.Address)
	// Redelegate 0.0005tkx to validator 2
	staking.MustRedelegate(t, val1.OperatorAddress, val2.OperatorAddress, "0.0005tkx", del.Address)
}

func TestUnbond(t *testing.T) {
	t.Parallel()

	// Create validator
	val := createValidator(t)
	// Create delegator account
	delName := testutil.GetName()
	defer testutil.PutName(delName)
	defer keys.MustDelete(t, delName)
	del := keys.MustAdd(t, delName)
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, del.Address, "1tkx")
	// Delegate 0.001tkx
	staking.MustDelegate(t, val.OperatorAddress, "0.001tkx", del.Address)
	// Unbond 0.0005tkx
	staking.MustUnbond(t, val.OperatorAddress, "0.0005tkx", del.Address)
}

func TestCancelUnbound(t *testing.T) {
	t.Parallel()

	// Create validator
	val := createValidator(t)
	// Create delegator account
	delName := testutil.GetName()
	defer testutil.PutName(delName)
	defer keys.MustDelete(t, delName)
	del := keys.MustAdd(t, delName)
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, del.Address, "1tkx")
	// Delegate 0.001tkx
	staking.MustDelegate(t, val.OperatorAddress, "0.001tkx", del.Address)
	// Unbond 0.0005tkx
	tx := staking.MustUnbond(t, val.OperatorAddress, "0.0005tkx", del.Address)
	// Cancel unbond 0.0002tkx
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.0002tkx", tx.Height.Int64(), del.Address)
}
