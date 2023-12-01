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
	defer keys.MustDeleteAccount(t, "alice") // Always delete account in case something went wrong
	alice := keys.MustCreateAccount(t, "alice")
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, alice, "1tkx")
	// Stake 0.001tkx
	return staking.MustCreateValidator(t, valPk, "0.001tkx", 0.1, 0.2, 0.001, 1, alice)
}

func TestCreateValidator(t *testing.T) {
	createValidator(t)
}

func TestDelegate(t *testing.T) {
	// Create validator
	val := createValidator(t)
	// Create delegator account
	defer keys.MustDeleteAccount(t, "bob") // Always delete account in case something went wrong
	bob := keys.MustCreateAccount(t, "bob")
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, bob, "1tkx")
	// Delegate 0.001tkx
	staking.MustDelegate(t, val.OperatorAddress, "0.001tkx", bob)
}

func TestRedelegate(t *testing.T) {
	// Create validator 1
	val1 := createValidator(t)
	// Create validator 2
	val2 := createValidator(t)
	// Create delegator account
	defer keys.MustDeleteAccount(t, "bob") // Always delete account in case something went wrong
	bob := keys.MustCreateAccount(t, "bob")
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, bob, "1tkx")
	// Delegate 0.001tkx to validator 1
	staking.MustDelegate(t, val1.OperatorAddress, "0.001tkx", bob)
	// Redelegate 0.0005tkx to validator 2
	staking.MustDelegate(t, val2.OperatorAddress, "0.0005tkx", bob)
}

func TestUnbond(t *testing.T) {
	// Create validator
	val := createValidator(t)
	// Create delegator account
	defer keys.MustDeleteAccount(t, "bob") // Always delete account in case something went wrong
	bob := keys.MustCreateAccount(t, "bob")
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, bob, "1tkx")
	// Delegate 0.001tkx
	staking.MustDelegate(t, val.OperatorAddress, "0.001tkx", bob)
	// Unbond 0.0005tkx
	staking.MustUnbond(t, val.OperatorAddress, "0.0005tkx", bob)
}

func TestCancelUnbound(t *testing.T) {
	// Create validator
	val := createValidator(t)
	// Create delegator account
	defer keys.MustDeleteAccount(t, "bob") // Always delete account in case something went wrong
	bob := keys.MustCreateAccount(t, "bob")
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, bob, "1tkx")
	// Delegate 0.001tkx
	staking.MustDelegate(t, val.OperatorAddress, "0.001tkx", bob)
	// Unbond 0.0005tkx
	tx := staking.MustUnbond(t, val.OperatorAddress, "0.0005tkx", bob)
	// Cancel unbond 0.0002tkx
	staking.MustCancelUnbound(t, val.OperatorAddress, "0.0002tkx", tx.Height.Int64(), bob)
}
