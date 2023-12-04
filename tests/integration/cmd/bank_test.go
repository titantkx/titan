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
	txcmd "github.com/tokenize-titan/titan/testutil/cmd/tx"
)

var faucetMtx sync.Mutex

func MustAcquireMoney(t testing.TB, address string, amount string) {
	faucet := keys.MustShowAddress(t, "faucet")
	faucetMtx.Lock()
	defer faucetMtx.Unlock()
	bank.MustSend(t, faucet, address, amount)
}

func TestTotalBalanceNotChanged(t *testing.T) {
	totalBal := testutil.MakeBigInt(0)
	accounts := auth.MustGetAccounts(t)

	var mtx sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for account := range accounts {
				bal := bank.MustGetBalance(t, account.GetAddress(), "utkx")
				mtx.Lock()
				totalBal = totalBal.Add(bal)
				mtx.Unlock()
			}
		}()
	}
	wg.Wait()

	require.Equal(t, testutil.MakeBigIntFromString("100000000000000000000000000"), totalBal)
}

func TestSend(t *testing.T) {
	t.Parallel()

	senderName := testutil.GetName()
	defer testutil.PutName(senderName)
	receiverName := testutil.GetName()
	defer testutil.PutName(receiverName)

	// Create sender account
	defer keys.MustDelete(t, senderName)
	sender := keys.MustAdd(t, senderName)
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, sender.Address, "1tkx")
	// Create receiver account
	defer keys.MustDelete(t, receiverName)
	receiver := keys.MustAdd(t, receiverName)
	// Send 0.5tkx
	bank.MustSend(t, sender.Address, receiver.Address, "0.5tkx")
}

func TestSendLowBalance(t *testing.T) {
	t.Parallel()

	senderName := testutil.GetName()
	defer testutil.PutName(senderName)
	receiverName := testutil.GetName()
	defer testutil.PutName(receiverName)

	// Create sender account
	defer keys.MustDelete(t, senderName)
	sender := keys.MustAdd(t, senderName)
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, sender.Address, "1tkx")
	// Create receiver account
	defer keys.MustDelete(t, receiverName)
	receiver := keys.MustAdd(t, receiverName)
	// Attempt to send 2tkx
	senderBalBefore := bank.MustGetBalance(t, sender.Address, "utkx")
	receiverBalBefore := bank.MustGetBalance(t, receiver.Address, "utkx")

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx, err := txcmd.ExecTx(ctx, "bank", "send", sender.Address, receiver.Address, "2tkx")

	require.Nil(t, tx)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient funds")

	senderBalAfter := bank.MustGetBalance(t, sender.Address, "utkx")
	receiverBalAfter := bank.MustGetBalance(t, receiver.Address, "utkx")

	require.Equal(t, senderBalBefore, senderBalAfter)
	require.Equal(t, receiverBalBefore, receiverBalAfter)
}

func TestMultiSend(t *testing.T) {
	t.Parallel()

	senderName := testutil.GetName()
	defer testutil.PutName(senderName)
	receiver1Name := testutil.GetName()
	defer testutil.PutName(receiver1Name)
	receiver2Name := testutil.GetName()
	defer testutil.PutName(receiver2Name)

	// Create sender account
	defer keys.MustDelete(t, senderName)
	sender := keys.MustAdd(t, senderName)
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, sender.Address, "1tkx")
	// Create receiver accounts
	defer keys.MustDelete(t, receiver1Name)
	receiver1 := keys.MustAdd(t, receiver1Name)
	defer keys.MustDelete(t, receiver2Name)
	receiver2 := keys.MustAdd(t, receiver2Name)
	// Send 0.3tkx to each receiver
	bank.MustMultiSend(t, sender.Address, "0.3tkx", receiver1.Address, receiver2.Address)
}

func TestMultiSendLowBalance(t *testing.T) {
	t.Parallel()

	senderName := testutil.GetName()
	defer testutil.PutName(senderName)
	receiver1Name := testutil.GetName()
	defer testutil.PutName(receiver1Name)
	receiver2Name := testutil.GetName()
	defer testutil.PutName(receiver2Name)

	// Create sender account
	defer keys.MustDelete(t, senderName)
	sender1 := keys.MustAdd(t, senderName)
	// Ask for 1tkx from the faucet
	MustAcquireMoney(t, sender1.Address, "1tkx")
	// Create receiver accounts
	defer keys.MustDelete(t, receiver1Name)
	receiver1 := keys.MustAdd(t, receiver1Name)
	defer keys.MustDelete(t, receiver2Name)
	receiver2 := keys.MustAdd(t, receiver2Name)
	// Attempt to send 0.6tkx to each receiver
	senderBalBefore := bank.MustGetBalance(t, sender1.Address, "utkx")
	receiver1BalBefore := bank.MustGetBalance(t, receiver1.Address, "utkx")
	receiver2BalBefore := bank.MustGetBalance(t, receiver2.Address, "utkx")

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx, err := txcmd.ExecTx(ctx, "bank", "multi-send", sender1.Address, receiver1.Address, receiver2.Address, "0.6tkx")

	require.Nil(t, tx)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient funds")

	senderBalAfter := bank.MustGetBalance(t, sender1.Address, "utkx")
	receiver1BalAfter := bank.MustGetBalance(t, receiver1.Address, "utkx")
	receiver2BalAfter := bank.MustGetBalance(t, receiver2.Address, "utkx")

	require.Equal(t, senderBalBefore, senderBalAfter)
	require.Equal(t, receiver1BalBefore, receiver1BalAfter)
	require.Equal(t, receiver2BalBefore, receiver2BalAfter)
}
