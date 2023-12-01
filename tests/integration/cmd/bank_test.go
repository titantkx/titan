package cmd_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd/bank"
	"github.com/tokenize-titan/titan/testutil/cmd/keys"
	txcmd "github.com/tokenize-titan/titan/testutil/cmd/tx"
)

func MustAcquireMoney(t testing.TB, address string, amount string) {
	faucet := keys.MustGetAddress(t, "faucet")
	bank.MustSend(t, faucet, address, amount)
}

func TestSend(t *testing.T) {
	// Create sender account
	defer keys.MustDeleteAccount(t, "alice") // Always delete account in case something went wrong
	alice := keys.MustCreateAccount(t, "alice")
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, alice, "1tkx")
	// Create receiver account
	defer keys.MustDeleteAccount(t, "bob") // Always delete account in case something went wrong
	bob := keys.MustCreateAccount(t, "bob")
	// Send 0.5tkx
	bank.MustSend(t, alice, bob, "0.5tkx")
}

func TestSendLowBalance(t *testing.T) {
	// Create sender account
	defer keys.MustDeleteAccount(t, "alice") // Always delete account in case something went wrong
	alice := keys.MustCreateAccount(t, "alice")
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, alice, "1tkx")
	// Create receiver account
	defer keys.MustDeleteAccount(t, "bob") // Always delete account in case something went wrong
	bob := keys.MustCreateAccount(t, "bob")
	// Attempt to send 2tkx
	aliceBalBefore := bank.MustGetBalance(t, alice, "utkx")
	bobBalBefore := bank.MustGetBalance(t, bob, "utkx")

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx, err := txcmd.ExecTx(ctx, "bank", "send", alice, bob, "2tkx")

	require.Nil(t, tx)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient funds")

	aliceBalAfter := bank.MustGetBalance(t, alice, "utkx")
	bobBalAfter := bank.MustGetBalance(t, bob, "utkx")

	require.Equal(t, aliceBalBefore, aliceBalAfter)
	require.Equal(t, bobBalBefore, bobBalAfter)
}

func TestMultiSend(t *testing.T) {
	// Create sender account
	defer keys.MustDeleteAccount(t, "alice") // Always delete account in case something went wrong
	alice := keys.MustCreateAccount(t, "alice")
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, alice, "1tkx")
	// Create receiver accounts
	defer keys.MustDeleteAccount(t, "bob") // Always delete account in case something went wrong
	bob := keys.MustCreateAccount(t, "bob")
	defer keys.MustDeleteAccount(t, "carol") // Always delete account in case something went wrong
	carol := keys.MustCreateAccount(t, "carol")
	// Send 0.3tkx to each receiver
	bank.MustMultiSend(t, alice, "0.3tkx", bob, carol)
}

func TestMultiSendLowBalance(t *testing.T) {
	// Create sender account
	defer keys.MustDeleteAccount(t, "alice") // Always delete account in case something went wrong
	alice := keys.MustCreateAccount(t, "alice")
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, alice, "1tkx")
	// Create receiver accounts
	defer keys.MustDeleteAccount(t, "bob") // Always delete account in case something went wrong
	bob := keys.MustCreateAccount(t, "bob")
	defer keys.MustDeleteAccount(t, "carol") // Always delete account in case something went wrong
	carol := keys.MustCreateAccount(t, "carol")
	// Attempt to send 0.6tkx to each receiver
	aliceBalBefore := bank.MustGetBalance(t, alice, "utkx")
	bobBalBefore := bank.MustGetBalance(t, bob, "utkx")
	carolBalBefore := bank.MustGetBalance(t, carol, "utkx")

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx, err := txcmd.ExecTx(ctx, "bank", "multi-send", alice, bob, carol, "0.6tkx")

	require.Nil(t, tx)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient funds")

	aliceBalAfter := bank.MustGetBalance(t, alice, "utkx")
	bobBalAfter := bank.MustGetBalance(t, bob, "utkx")
	carolBalAfter := bank.MustGetBalance(t, carol, "utkx")

	require.Equal(t, aliceBalBefore, aliceBalAfter)
	require.Equal(t, bobBalBefore, bobBalAfter)
	require.Equal(t, carolBalBefore, carolBalAfter)
}
