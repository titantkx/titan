package cmd_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd/auth"
	"github.com/tokenize-titan/titan/testutil/cmd/bank"
	"github.com/tokenize-titan/titan/testutil/cmd/keys"
	"github.com/tokenize-titan/titan/testutil/cmd/status"
	txcmd "github.com/tokenize-titan/titan/testutil/cmd/tx"
)

var faucetMtx sync.Mutex

func MustAcquireMoney(t testing.TB, address string, amount string) {
	faucet := keys.MustShowAddress(t, "faucet")
	faucetMtx.Lock()
	defer faucetMtx.Unlock()
	bank.MustSend(t, faucet, address, amount)
}

func MustCreateAccount(t testing.TB, balance string) keys.Key {
	key := MustAddKey(t)
	if balance != "" {
		MustAcquireMoney(t, key.Address, balance)
	}
	return key
}

func MustGetTotalBalance(t testing.TB, height int64) testutil.BigInt {
	if height <= 0 {
		height = status.MustGetStatus(t).SyncInfo.LatestBlockHeight.Int64()
	}

	accounts := auth.MustGetAccounts(t, height)
	totalBal := testutil.MakeBigInt(0)

	var mtx sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for account := range accounts {
				bal := bank.MustGetBalance(t, account.GetAddress(), "utkx", height)
				mtx.Lock()
				totalBal = totalBal.Add(bal)
				mtx.Unlock()
			}
		}()
	}
	wg.Wait()

	return totalBal
}

func TestTotalBalanceNotChanged(t *testing.T) {
	t.Parallel()

	totalBal := MustGetTotalBalance(t, 0)

	require.Equal(t, testutil.MakeBigIntFromString("100000000000000000000000000"), totalBal)
}

func TestSend(t *testing.T) {
	t.Parallel()

	sender := MustCreateAccount(t, "1tkx")
	receiver := MustCreateAccount(t, "")

	bank.MustSend(t, sender.Address, receiver.Address, "0.5tkx")
}

func TestSendLowBalance(t *testing.T) {
	t.Parallel()

	sender := MustCreateAccount(t, "1tkx")
	receiver := MustCreateAccount(t, "")

	// Attempt to send 2tkx
	senderBalBefore := bank.MustGetBalance(t, sender.Address, "utkx", 0)
	receiverBalBefore := bank.MustGetBalance(t, receiver.Address, "utkx", 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx, err := txcmd.ExecTx(ctx, "bank", "send", sender.Address, receiver.Address, "2tkx")

	require.Nil(t, tx)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient funds")

	senderBalAfter := bank.MustGetBalance(t, sender.Address, "utkx", 0)
	receiverBalAfter := bank.MustGetBalance(t, receiver.Address, "utkx", 0)

	require.Equal(t, senderBalBefore, senderBalAfter)
	require.Equal(t, receiverBalBefore, receiverBalAfter)
}

func TestMultiSend(t *testing.T) {
	t.Parallel()

	sender := MustCreateAccount(t, "1tkx")
	receiver1 := MustCreateAccount(t, "")
	receiver2 := MustCreateAccount(t, "")

	bank.MustMultiSend(t, sender.Address, "0.3tkx", receiver1.Address, receiver2.Address)
}

func TestMultiSendLowBalance(t *testing.T) {
	t.Parallel()

	sender := MustCreateAccount(t, "1tkx")
	receiver1 := MustCreateAccount(t, "")
	receiver2 := MustCreateAccount(t, "")

	// Attempt to send 0.6tkx to each receiver
	senderBalBefore := bank.MustGetBalance(t, sender.Address, "utkx", 0)
	receiver1BalBefore := bank.MustGetBalance(t, receiver1.Address, "utkx", 0)
	receiver2BalBefore := bank.MustGetBalance(t, receiver2.Address, "utkx", 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx, err := txcmd.ExecTx(ctx, "bank", "multi-send", sender.Address, receiver1.Address, receiver2.Address, "0.6tkx")

	require.Nil(t, tx)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient funds")

	senderBalAfter := bank.MustGetBalance(t, sender.Address, "utkx", 0)
	receiver1BalAfter := bank.MustGetBalance(t, receiver1.Address, "utkx", 0)
	receiver2BalAfter := bank.MustGetBalance(t, receiver2.Address, "utkx", 0)

	require.Equal(t, senderBalBefore, senderBalAfter)
	require.Equal(t, receiver1BalBefore, receiver1BalAfter)
	require.Equal(t, receiver2BalBefore, receiver2BalAfter)
}
