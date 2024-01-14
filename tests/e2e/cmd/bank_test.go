package cmd_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tokenize-titan/titan/utils"

	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
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
				bal := bank.MustGetBalance(t, account.GetAddress(), utils.BaseDenom, height)
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

	sender := MustCreateAccount(t, "1"+utils.DisplayDenom)
	receiver := MustCreateAccount(t, "")

	bank.MustSend(t, sender.Address, receiver.Address, "0.5"+utils.DisplayDenom)
}

func TestSendLowBalance(t *testing.T) {
	t.Parallel()

	sender := MustCreateAccount(t, "1"+utils.DisplayDenom)
	receiver := MustCreateAccount(t, "")

	// Attempt to send 2tkx
	senderBalBefore := bank.MustGetBalance(t, sender.Address, utils.BaseDenom, 0)
	receiverBalBefore := bank.MustGetBalance(t, receiver.Address, utils.BaseDenom, 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx, err := txcmd.ExecTx(ctx, "bank", "send", sender.Address, receiver.Address, "2"+utils.DisplayDenom)

	require.Nil(t, tx)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient funds")

	senderBalAfter := bank.MustGetBalance(t, sender.Address, utils.BaseDenom, 0)
	receiverBalAfter := bank.MustGetBalance(t, receiver.Address, utils.BaseDenom, 0)

	require.Equal(t, senderBalBefore, senderBalAfter)
	require.Equal(t, receiverBalBefore, receiverBalAfter)
}

func TestMultiSend(t *testing.T) {
	t.Parallel()

	sender := MustCreateAccount(t, "1"+utils.DisplayDenom)
	receiver1 := MustCreateAccount(t, "")
	receiver2 := MustCreateAccount(t, "")

	bank.MustMultiSend(t, sender.Address, "0.3"+utils.DisplayDenom, receiver1.Address, receiver2.Address)
}

func TestMultiSendLowBalance(t *testing.T) {
	t.Parallel()

	sender := MustCreateAccount(t, "1"+utils.DisplayDenom)
	receiver1 := MustCreateAccount(t, "")
	receiver2 := MustCreateAccount(t, "")

	// Attempt to send 0.6tkx to each receiver
	senderBalBefore := bank.MustGetBalance(t, sender.Address, utils.BaseDenom, 0)
	receiver1BalBefore := bank.MustGetBalance(t, receiver1.Address, utils.BaseDenom, 0)
	receiver2BalBefore := bank.MustGetBalance(t, receiver2.Address, utils.BaseDenom, 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx, err := txcmd.ExecTx(ctx, "bank", "multi-send", sender.Address, receiver1.Address, receiver2.Address, "0.6"+utils.DisplayDenom)

	require.Nil(t, tx)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient funds")

	senderBalAfter := bank.MustGetBalance(t, sender.Address, utils.BaseDenom, 0)
	receiver1BalAfter := bank.MustGetBalance(t, receiver1.Address, utils.BaseDenom, 0)
	receiver2BalAfter := bank.MustGetBalance(t, receiver2.Address, utils.BaseDenom, 0)

	require.Equal(t, senderBalBefore, senderBalAfter)
	require.Equal(t, receiver1BalBefore, receiver1BalAfter)
	require.Equal(t, receiver2BalBefore, receiver2BalAfter)
}

func TestMultisigSend(t *testing.T) {
	t.Parallel()

	sender1 := MustCreateAccount(t, "")
	sender2 := MustCreateAccount(t, "")
	receiver := MustCreateAccount(t, "")
	multisigAccount := MustAddMultisigKey(t, 2, sender1.Name, sender2.Name)

	MustAcquireMoney(t, multisigAccount.Address, "1"+utils.DisplayDenom)

	unsingedTx := testutil.MustCreateTemp(t, "unsigned_tx_*.json")
	unsignedTxContent := txcmd.MustGenerateTx(t, "bank", "send", multisigAccount.Name, receiver.Address, "0.5"+utils.DisplayDenom)
	_, err := unsingedTx.Write(unsignedTxContent)
	require.NoError(t, err)

	signedTx1 := testutil.MustCreateTemp(t, "singed_tx_1_*.json")
	signedTx1Content := cmd.MustExec(t, "titand", "tx", "sign", "--from="+sender1.Name, "--multisig="+multisigAccount.Address, unsingedTx.Name())
	require.NotEmpty(t, signedTx1Content)
	_, err = signedTx1.Write(signedTx1Content)
	require.NoError(t, err)

	signedTx2 := testutil.MustCreateTemp(t, "singed_tx_2_*.json")
	signedTx2Content := cmd.MustExec(t, "titand", "tx", "sign", "--from="+sender2.Name, "--multisig="+multisigAccount.Address, unsingedTx.Name())
	require.NotEmpty(t, signedTx2Content)
	_, err = signedTx2.Write(signedTx2Content)
	require.NoError(t, err)

	signedTx := testutil.MustCreateTemp(t, "singed_tx_*.json")
	signedTxContent := cmd.MustExec(t, "titand", "tx", "multisign", unsingedTx.Name(), multisigAccount.Name, signedTx1.Name(), signedTx2.Name())
	require.NotEmpty(t, signedTxContent)
	_, err = signedTx.Write(signedTxContent)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()
	tx := txcmd.MustBroadcastTx(t, ctx, signedTx.Name())

	coinSpent := tx.MustGetDeductFeeAmount(t)
	senderBal := bank.MustGetBalance(t, multisigAccount.Address, utils.BaseDenom, 0)
	receiverBal := bank.MustGetBalance(t, receiver.Address, utils.BaseDenom, 0)

	require.Equal(t, testutil.MakeBigIntFromString("500000000000000000").Sub(coinSpent), senderBal)
	require.Equal(t, testutil.MakeBigIntFromString("500000000000000000"), receiverBal)
}
