package ibc_hook_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	interchaintest "github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibchookskeeper "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/keeper"

	"github.com/titantkx/titan/tests/interchain"
)

type WasmMetadata struct {
	Contract string                            `json:"contract"`
	Msg      map[string]map[string]interface{} `json:"msg"`
}

type PacketMetadata struct {
	Wasm *WasmMetadata `json:"wasm"`
}

type QueryParams struct {
	Address string `json:"addr"`
}

type CounterQueryGetCount struct {
	GetCount *QueryParams `json:"get_count"`
}
type CounterQueryGetTotalFunds struct {
	GetTotalFunds *QueryParams `json:"get_total_funds"`
}

type CounterResponseGetCount struct {
	Data struct {
		Count int `json:"count"`
	} `json:"data"`
}

type CounterResponseGetTotalFunds struct {
	Data struct {
		TotalFunds []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"total_funds"`
	} `json:"data"`
}

func TestIbcHook(t *testing.T) {
	// Set an environment variable `KEEP_CONTAINERS` if it's not already set
	if os.Getenv("KEEP_CONTAINERS") == "" {
		err := os.Setenv("KEEP_CONTAINERS", "")
		if err != nil {
			t.Logf("Error setting environment variable: %s", err)
			return
		}
	}

	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	var (
		ctx                = context.Background()
		client, network    = interchaintest.DockerSetup(t)
		rep                = testreporter.NewNopReporter()
		eRep               = rep.RelayerExecReporter(t)
		chainIDT, chainID1 = "titan_18887-1", "cosmoshub-1"
	)

	numValidators := 1 // Defines how many validators should be used in each network.
	numFullNodes := 0  // Defines how many additional full nodes should be used in each network.

	// Here we define our ChainFactory by instantiating a new instance of the BuiltinChainFactory exposed in interchaintest.
	// We use the ChainSpec type to fully describe which chains we want to use in our tests.
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		interchain.TitanChainSpec(ctx, chainIDT, numValidators, numFullNodes),
		{
			Name:          "gaia",
			ChainName:     chainID1,
			Version:       "v14.1.0",
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
			ChainConfig: ibc.ChainConfig{
				ChainID: chainID1,
			},
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)
	chainT, chain1 := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
	).Build(t, client, network)

	const pathT1 = "t1"

	// build the interchain
	ic := interchaintest.NewInterchain().
		AddChain(chainT).
		AddChain(chain1).
		AddRelayer(r, "relayer").
		AddLink(interchaintest.InterchainLink{
			Chain1:  chainT,
			Chain2:  chain1,
			Relayer: r,
			Path:    pathT1,
		})
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),

		SkipPathCreation: false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	t1Chan, err := ibc.GetTransferChannel(ctx, r, eRep, chainIDT, chainID1)
	require.NoError(t, err)
	c1tChan := t1Chan.Counterparty

	// Start the relayer on all paths
	err = r.StartRelayer(ctx, eRep, pathT1)
	require.NoError(t, err)
	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	c1tDenom := transfertypes.GetPrefixedDenom(t1Chan.PortID, t1Chan.ChannelID, chain1.Config().Denom)
	c1tDenomTrace := transfertypes.ParseDenomTrace(c1tDenom)
	c1tEscrowAccount := sdk.MustBech32ifyAddressBytes(chain1.Config().Bech32Prefix, transfertypes.GetEscrowAddress(c1tChan.PortID, c1tChan.ChannelID))

	//////////////////////////

	// Create and Fund User Wallets
	userFunds := sdkmath.NewInt(10_000) // *10^decimals
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chainT, chain1)
	userT, user1 := users[0], users[1]

	require.NoError(t, testutil.WaitForBlocks(ctx, 1, chainT, chain1))

	// Store counter.wasm contract
	counterCodeId, err := chainT.StoreContract(
		ctx, userT.KeyName(), "wasm/counter/artifacts/counter.wasm")
	require.NoError(t, err)

	// Instantiate counter.wasm contract
	counterContractAddr, err := chainT.InstantiateContract(
		ctx, userT.KeyName(), counterCodeId, `{"count":0}`, true)
	require.NoError(t, err)

	require.NoError(t, testutil.WaitForBlocks(ctx, 1, chainT, chain1))

	t.Run("trigger contract logic 1->t", func(t *testing.T) {
		amountToSend := sdkmath.NewInt(10)

		transfer := ibc.WalletAmount{
			Address: counterContractAddr,
			Denom:   chain1.Config().Denom,
			Amount:  amountToSend,
		}

		metadata := &PacketMetadata{
			Wasm: &WasmMetadata{
				Contract: counterContractAddr,
				Msg: map[string]map[string]interface{}{
					"increment": {},
				},
			},
		}
		memo, err := json.Marshal(metadata)
		require.NoError(t, err)

		timesToSend := 2
		for i := 0; i < timesToSend; i++ {
			height, err := chain1.Height(ctx)
			require.NoError(t, err)
			transferTx, err := chain1.SendIBCTransfer(ctx, c1tChan.ChannelID, user1.KeyName(), transfer, ibc.TransferOptions{Memo: string(memo)})
			require.NoError(t, err)
			_, err = testutil.PollForAck(ctx, chain1, height, height+30, transferTx.Packet)
			require.NoError(t, err)
			err = testutil.WaitForBlocks(ctx, 1, chain1, chainT)
			require.NoError(t, err)
		}

		t.Logf("sender: %s", user1.FormattedAddress())
		t.Logf("receiver: %s", userT.FormattedAddress())
		t.Logf("channel: %s > %s", c1tChan.ChannelID, t1Chan.ChannelID)

		// get the derived account to check the count
		senderBech32, err := ibchookskeeper.DeriveIntermediateSender(
			t1Chan.ChannelID,
			user1.FormattedAddress(),
			chainT.Config().Bech32Prefix,
		)
		t.Logf("senderBech32: %s", senderBech32)
		require.NoError(t, err)

		{
			query := fmt.Sprintf(`{"get_count": {"addr": "%s"}}`, senderBech32)
			stdout, _, err := chainT.GetNode().ExecQuery(ctx, "wasm", "contract-state", "smart", counterContractAddr, query)
			require.NoError(t, err)
			results := &CounterResponseGetCount{}
			err = json.Unmarshal(stdout, results)
			require.NoError(t, err)
			require.Equal(t, 1, results.Data.Count)
			t.Logf("stdout: %s", stdout)
		}

		{
			query := fmt.Sprintf(`{"get_total_funds": {"addr": "%s"}}`, senderBech32)
			stdout, _, err := chainT.GetNode().ExecQuery(ctx, "wasm", "contract-state", "smart", counterContractAddr, query)
			require.NoError(t, err)
			results := &CounterResponseGetTotalFunds{}
			err = json.Unmarshal(stdout, results)
			require.NoError(t, err)
			require.Equal(t, timesToSend-1, len(results.Data.TotalFunds))
			require.Equal(t, c1tDenomTrace.IBCDenom(), results.Data.TotalFunds[0].Denom)
			require.Equal(t, amountToSend.MulRaw(int64(timesToSend)).String(), results.Data.TotalFunds[0].Amount)
			t.Logf("stdout: %s", stdout)
		}
	})

	t.Run("trigger fail contract logic 1->t", func(t *testing.T) {
		amountToSend := sdkmath.NewInt(10)

		transfer := ibc.WalletAmount{
			Address: counterContractAddr,
			Denom:   chain1.Config().Denom,
			Amount:  amountToSend,
		}

		metadata := &PacketMetadata{
			Wasm: &WasmMetadata{
				Contract: counterContractAddr,
				Msg: map[string]map[string]interface{}{
					"increment_not_existed": {},
				},
			},
		}
		memo, err := json.Marshal(metadata)
		require.NoError(t, err)

		chain1OldBalance, err := chain1.GetBalance(ctx, user1.FormattedAddress(), chain1.Config().Denom)
		require.NoError(t, err)
		c1tEscrowOldBalance, err := chain1.GetBalance(ctx, c1tEscrowAccount, chain1.Config().Denom)
		require.NoError(t, err)
		chainTContractOldBalance, err := chainT.GetBalance(ctx, counterContractAddr, c1tDenomTrace.IBCDenom())
		require.NoError(t, err)

		height, err := chain1.Height(ctx)
		require.NoError(t, err)
		transferTx, err := chain1.SendIBCTransfer(ctx, c1tChan.ChannelID, user1.KeyName(), transfer, ibc.TransferOptions{Memo: string(memo)})
		require.NoError(t, err)
		pkack, err := testutil.PollForAck(ctx, chain1, height, height+30, transferTx.Packet)
		require.NoError(t, err)
		err = testutil.WaitForBlocks(ctx, 1, chain1, chainT)
		require.NoError(t, err)
		t.Logf("ack: %s", string(pkack.Acknowledgement))
		var ack map[string]string // This can't be unmarshalled to Acknowledgement because it's fetched from the events
		err = json.Unmarshal(pkack.Acknowledgement, &ack)
		require.NoError(t, err)
		require.Contains(t, ack, "error")

		chain1Balance, err := chain1.GetBalance(ctx, user1.FormattedAddress(), chain1.Config().Denom)
		require.NoError(t, err)
		c1tEscrowBalance, err := chain1.GetBalance(ctx, c1tEscrowAccount, chain1.Config().Denom)
		require.NoError(t, err)
		chainTContractBalance, err := chainT.GetBalance(ctx, counterContractAddr, c1tDenomTrace.IBCDenom())
		require.NoError(t, err)

		require.Equal(t, chain1OldBalance.Sub(transferTx.Fee), chain1Balance)
		require.Equal(t, c1tEscrowOldBalance, c1tEscrowBalance)
		require.Equal(t, chainTContractOldBalance, chainTContractBalance)
	})
}
