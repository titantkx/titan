package interchain_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	interchaintest "github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v7/relayer"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/titantkx/titan/tests/interchain"
)

type ForwardMetadata struct {
	Receiver       string        `json:"receiver"`
	Port           string        `json:"port"`
	Channel        string        `json:"channel"`
	Timeout        time.Duration `json:"timeout"`
	Retries        *uint8        `json:"retries,omitempty"`
	Next           *string       `json:"next,omitempty"`
	RefundSequence *uint64       `json:"refund_sequence,omitempty"`
}

type PacketMetadata struct {
	Forward *ForwardMetadata `json:"forward"`
}

// TestPacketForwardMiddlewareRouter ensures the PFM module is set up properly and works as expected.
func TestPacketForwardMiddlewareRouter(t *testing.T) {
	// Set an environment variable
	err := os.Setenv("KEEP_CONTAINERS", "")
	if err != nil {
		t.Logf("Error setting environment variable: %s", err)
		return
	}

	if testing.Short() {
		t.Skip()
	}

	numValidators, numFullNodes := 1, 0

	var (
		ctx                          = context.Background()
		client, network              = interchaintest.DockerSetup(t)
		rep                          = testreporter.NewNopReporter()
		eRep                         = rep.RelayerExecReporter(t)
		chainIDT, chainID1, chainID2 = "titan_18887-1", "cosmoshub-1", "cosmoshub-2"
		chainT, chain1, chain2       *cosmos.CosmosChain
	)

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
		{
			Name:          "gaia",
			ChainName:     chainID2,
			Version:       "v14.1.0",
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
			ChainConfig: ibc.ChainConfig{
				ChainID: chainID2,
			},
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chainT, chain1, chain2 = chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain), chains[2].(*cosmos.CosmosChain)

	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
		interchaintestrelayer.CustomDockerImage(interchain.IBCRelayerImage, interchain.IBCRelayerVersion, "1000:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "100", "--log-level", "debug"), // , "--override"
	).Build(t, client, network)

	const pathT1 = "t1"
	const pathT2 = "t2"

	ic := interchaintest.NewInterchain().
		AddChain(chainT).
		AddChain(chain1).
		AddChain(chain2).
		AddRelayer(r, "relayer").
		AddLink(interchaintest.InterchainLink{
			Chain1:  chainT,
			Chain2:  chain1,
			Relayer: r,
			Path:    pathT1,
		}).
		AddLink(interchaintest.InterchainLink{
			Chain1:  chainT,
			Chain2:  chain2,
			Relayer: r,
			Path:    pathT2,
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

	userFunds := sdkmath.NewInt(10_000) // *10^decimals
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chainT, chain1, chain2)

	t1Chan, err := ibc.GetTransferChannel(ctx, r, eRep, chainIDT, chainID1)
	require.NoError(t, err)
	c1tChan := t1Chan.Counterparty

	t2Chan, err := ibc.GetTransferChannel(ctx, r, eRep, chainIDT, chainID2)
	require.NoError(t, err)
	c2tChan := t2Chan.Counterparty

	// Start the relayer on all paths
	err = r.StartRelayer(ctx, eRep, pathT1, pathT2)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	// Get original account balances
	userT, user1, user2 := users[0], users[1], users[2]
	t.Logf("userT: %s, user1: %s, user2: %s", userT.FormattedAddress(), user1.FormattedAddress(), user2.FormattedAddress())

	transferAmount := sdkmath.NewInt(1).Mul(chain1.GetDecimalPow())

	// ew will transfer tokens from Gaia1 to Titan to Gaia2

	// Compose the prefixed denoms and ibc denom for asserting balances
	firstHopDenom := transfertypes.GetPrefixedDenom(t1Chan.PortID, t1Chan.ChannelID, chain1.Config().Denom)
	secondHopDenom := transfertypes.GetPrefixedDenom(c2tChan.PortID, c2tChan.ChannelID, firstHopDenom)

	firstHopDenomTrace := transfertypes.ParseDenomTrace(firstHopDenom)
	secondHopDenomTrace := transfertypes.ParseDenomTrace(secondHopDenom)

	firstHopIBCDenom := firstHopDenomTrace.IBCDenom()
	secondHopIBCDenom := secondHopDenomTrace.IBCDenom()

	firstHopEscrowAccount := sdk.MustBech32ifyAddressBytes(chain1.Config().Bech32Prefix, transfertypes.GetEscrowAddress(c1tChan.PortID, c1tChan.ChannelID))
	secondHopEscrowAccount := sdk.MustBech32ifyAddressBytes(chainT.Config().Bech32Prefix, transfertypes.GetEscrowAddress(t2Chan.PortID, t2Chan.ChannelID))

	getBalances := func() (chain1Balance sdkmath.Int, chainTBalance sdkmath.Int, chain2Balance sdkmath.Int, firstHopEscrowBalance sdkmath.Int, secondHopEscrowBalance sdkmath.Int) {
		chain1Balance, err := chain1.GetBalance(ctx, user1.FormattedAddress(), chain1.Config().Denom)
		require.NoError(t, err)
		chainTBalance, err = chainT.GetBalance(ctx, userT.FormattedAddress(), firstHopIBCDenom)
		require.NoError(t, err)
		chain2Balance, err = chain2.GetBalance(ctx, user2.FormattedAddress(), secondHopIBCDenom)
		require.NoError(t, err)

		firstHopEscrowBalance, err = chain1.GetBalance(ctx, firstHopEscrowAccount, chain1.Config().Denom)
		require.NoError(t, err)
		secondHopEscrowBalance, err = chainT.GetBalance(ctx, secondHopEscrowAccount, firstHopIBCDenom)
		require.NoError(t, err)
		return
	}

	t.Run("multi-hop 1->t->2", func(t *testing.T) {
		// Send packet from Chain 1 to Titan to Chain 2

		transfer := ibc.WalletAmount{
			Address: userT.FormattedAddress(),
			Denom:   chain1.Config().Denom,
			Amount:  transferAmount,
		}

		firstHopMetadata := &PacketMetadata{
			Forward: &ForwardMetadata{
				Receiver: user2.FormattedAddress(),
				Channel:  t2Chan.ChannelID,
				Port:     t2Chan.PortID,
			},
		}

		memo, err := json.Marshal(firstHopMetadata)
		require.NoError(t, err)

		chain1Height, err := chain1.Height(ctx)
		require.NoError(t, err)

		transferTx, err := chain1.SendIBCTransfer(ctx, c1tChan.ChannelID, user1.KeyName(), transfer, ibc.TransferOptions{Memo: string(memo)})
		require.NoError(t, err)
		_, err = testutil.PollForAck(ctx, chain1, chain1Height, chain1Height+30, transferTx.Packet)
		require.NoError(t, err)
		err = testutil.WaitForBlocks(ctx, 1, chain1)
		require.NoError(t, err)

		chain1Balance, chainTBalance, chain2Balance, firstHopEscrowBalance, secondHopEscrowBalance := getBalances()

		require.Equal(t, userFunds.Mul(chain1.GetDecimalPow()).Sub(transferAmount).Sub(transferTx.Fee), chain1Balance)
		require.Equal(t, transferAmount, firstHopEscrowBalance)

		require.Equal(t, sdkmath.NewInt(0), chainTBalance)
		require.Equal(t, transferAmount, secondHopEscrowBalance)

		require.Equal(t, transferAmount, chain2Balance)
	})

	t.Run("multi-hop unwind 2->t->1", func(t *testing.T) {
		// Send packet from Chain 2 to Titan to Chain 1

		transfer := ibc.WalletAmount{
			Address: userT.FormattedAddress(),
			Denom:   secondHopIBCDenom,
			Amount:  transferAmount,
		}

		firstHopMetadata := &PacketMetadata{
			Forward: &ForwardMetadata{
				Receiver: user1.FormattedAddress(),
				Channel:  t1Chan.ChannelID,
				Port:     t1Chan.PortID,
			},
		}

		memo, err := json.Marshal(firstHopMetadata)
		require.NoError(t, err)

		chain2Height, err := chain2.Height(ctx)
		require.NoError(t, err)

		chain1OldBalance, _, _, _, _ := getBalances()

		transferTx, err := chain2.SendIBCTransfer(ctx, c2tChan.ChannelID, user2.KeyName(), transfer, ibc.TransferOptions{Memo: string(memo)})
		require.NoError(t, err)
		_, err = testutil.PollForAck(ctx, chain2, chain2Height, chain2Height+30, transferTx.Packet)
		require.NoError(t, err)
		err = testutil.WaitForBlocks(ctx, 1, chain2)
		require.NoError(t, err)

		chain1Balance, chainTBalance, chain2Balance, firstHopEscrowBalance, secondHopEscrowBalance := getBalances()

		require.Equal(t, sdkmath.NewInt(0), chain2Balance)

		require.Equal(t, sdkmath.NewInt(0), secondHopEscrowBalance)
		require.Equal(t, sdkmath.NewInt(0), chainTBalance)

		require.Equal(t, sdkmath.NewInt(0), firstHopEscrowBalance)
		require.Equal(t, chain1OldBalance.Add(transferAmount), chain1Balance)
	})

	t.Run("forward ack error refund", func(t *testing.T) {
		// Send a malformed packet with invalid receiver address from Chain 1->Chain T->Chain 2
		// This should succeed in the first hop and fail to make the second hop; funds should then be refunded to Chain A.
		transfer := ibc.WalletAmount{
			Address: userT.FormattedAddress(),
			Denom:   chain1.Config().Denom,
			Amount:  transferAmount,
		}

		metadata := &PacketMetadata{
			Forward: &ForwardMetadata{
				Receiver: "xyz1t8eh66t2w5k67kwurmn5gqhtq6d2ja0vp7jmmq", // malformed receiver address on Chain C
				Channel:  t2Chan.ChannelID,
				Port:     t2Chan.PortID,
			},
		}

		memo, err := json.Marshal(metadata)
		require.NoError(t, err)

		chain1Height, err := chain1.Height(ctx)
		require.NoError(t, err)

		chain1OldBalance, chainTOldBalance, chain2OldBalance, _, _ := getBalances()

		transferTx, err := chain1.SendIBCTransfer(ctx, c1tChan.ChannelID, user1.KeyName(), transfer, ibc.TransferOptions{Memo: string(memo)})
		require.NoError(t, err)
		_, err = testutil.PollForAck(ctx, chain1, chain1Height, chain1Height+25, transferTx.Packet)
		require.NoError(t, err)
		err = testutil.WaitForBlocks(ctx, 1, chain1)
		require.NoError(t, err)

		chain1Balance, chainTBalance, chain2Balance, firstHopEscrowBalance, secondHopEscrowBalance := getBalances()

		require.Equal(t, chain1OldBalance.Sub(transferTx.Fee), chain1Balance)
		require.Equal(t, sdkmath.NewInt(0), firstHopEscrowBalance)

		require.Equal(t, chainTOldBalance, chainTBalance)
		require.Equal(t, sdkmath.NewInt(0), secondHopEscrowBalance)

		require.Equal(t, chain2OldBalance, chain2Balance)
	})

	t.Run("forward timeout refund", func(t *testing.T) {
		// Send packet from Chain 1->Chain T->Chain 2 with the timeout so low for T->2 transfer that it can not make it from B to C, which should result in a refund from T to 1 after two retries.
		transfer := ibc.WalletAmount{
			Address: userT.FormattedAddress(),
			Denom:   chain1.Config().Denom,
			Amount:  transferAmount,
		}

		retries := uint8(2)
		metadata := &PacketMetadata{
			Forward: &ForwardMetadata{
				Receiver: user2.FormattedAddress(),
				Channel:  t2Chan.ChannelID,
				Port:     t2Chan.PortID,
				Retries:  &retries,
				Timeout:  1 * time.Second,
			},
		}

		memo, err := json.Marshal(metadata)
		require.NoError(t, err)

		chain1Height, err := chain1.Height(ctx)
		require.NoError(t, err)

		chain1OldBalance, chainTOldBalance, chain2OldBalance, _, _ := getBalances()

		transferTx, err := chain1.SendIBCTransfer(ctx, c1tChan.ChannelID, user1.KeyName(), transfer, ibc.TransferOptions{Memo: string(memo)})
		require.NoError(t, err)
		_, err = testutil.PollForAck(ctx, chain1, chain1Height, chain1Height+25, transferTx.Packet)
		require.NoError(t, err)
		err = testutil.WaitForBlocks(ctx, 1, chain1)
		require.NoError(t, err)

		chain1Balance, chainTBalance, chain2Balance, firstHopEscrowBalance, secondHopEscrowBalance := getBalances()

		require.Equal(t, chain1OldBalance.Sub(transferTx.Fee), chain1Balance)
		require.Equal(t, sdkmath.NewInt(0), firstHopEscrowBalance)

		require.Equal(t, chainTOldBalance, chainTBalance)
		require.Equal(t, sdkmath.NewInt(0), secondHopEscrowBalance)

		require.Equal(t, chain2OldBalance, chain2Balance)
	})

	t.Run("multi-hop through native chain ack error refund", func(t *testing.T) {
		// send normal IBC transfer from T->1 to get funds in IBC denom, then do multihop err 1->T(native)->2
		// Compose the prefixed denoms and ibc denom for asserting balances
		t1Denom := transfertypes.GetPrefixedDenom(c1tChan.PortID, c1tChan.ChannelID, chainT.Config().Denom)
		t2Denom := transfertypes.GetPrefixedDenom(c2tChan.PortID, c2tChan.ChannelID, chainT.Config().Denom)

		t1DenomTrace := transfertypes.ParseDenomTrace(t1Denom)
		t2DenomTrace := transfertypes.ParseDenomTrace(t2Denom)

		t1IBCDenom := t1DenomTrace.IBCDenom()
		t2IBCDenom := t2DenomTrace.IBCDenom()

		transfer := ibc.WalletAmount{
			Address: user1.FormattedAddress(),
			Denom:   chainT.Config().Denom,
			Amount:  transferAmount,
		}

		chainTOld1Balance, err := chainT.GetBalance(ctx, userT.FormattedAddress(), chainT.Config().Denom)
		require.NoError(t, err)

		chainTHeight, err := chainT.Height(ctx)
		require.NoError(t, err)
		transferTx, err := chainT.SendIBCTransfer(ctx, t1Chan.ChannelID, userT.KeyName(), transfer, ibc.TransferOptions{})
		require.NoError(t, err)
		_, err = testutil.PollForAck(ctx, chainT, chainTHeight, chainTHeight+10, transferTx.Packet)
		require.NoError(t, err)
		err = testutil.WaitForBlocks(ctx, 1, chainT)
		require.NoError(t, err)

		chain1OldBalance, err := chain1.GetBalance(ctx, user1.FormattedAddress(), t1IBCDenom)
		require.NoError(t, err)
		chainTOld2Balance, err := chainT.GetBalance(ctx, userT.FormattedAddress(), chainT.Config().Denom)
		require.NoError(t, err)

		t1EscrowAccount := sdk.MustBech32ifyAddressBytes(chainT.Config().Bech32Prefix, transfertypes.GetEscrowAddress(t1Chan.PortID, t1Chan.ChannelID))
		t.Logf("t1EscrowAccount: %s", t1EscrowAccount)
		t1EscrowBalance, err := chainT.GetBalance(ctx, t1EscrowAccount, chainT.Config().Denom)
		require.NoError(t, err)

		require.Equal(t, transferAmount, chain1OldBalance)
		require.Equal(t, chainTOld1Balance.Sub(transferAmount).Sub(transferTx.Fee), chainTOld2Balance)
		require.Equal(t, transferAmount, t1EscrowBalance)

		// Send a malformed packet with invalid receiver address from Chain 1->Chain T->Chain 2
		// This should succeed in the first hop , then fail to make the second hop.
		// Funds should be refunded to Chain T and then to Chain 1 via acknowledgements with errors.
		transfer = ibc.WalletAmount{
			Address: userT.FormattedAddress(),
			Denom:   t1IBCDenom,
			Amount:  transferAmount,
		}

		firstHopMetadata := &PacketMetadata{
			Forward: &ForwardMetadata{
				Receiver: "xyz1t8eh66t2w5k67kwurmn5gqhtq6d2ja0vp7jmmq", // malformed receiver address on chain D
				Channel:  t2Chan.ChannelID,
				Port:     t2Chan.PortID,
			},
		}
		memo, err := json.Marshal(firstHopMetadata)
		require.NoError(t, err)

		chain2OldBalance, err := chain2.GetBalance(ctx, user2.FormattedAddress(), t2IBCDenom)
		require.NoError(t, err)

		chain1Height, err := chain1.Height(ctx)
		require.NoError(t, err)
		transferTx, err = chain1.SendIBCTransfer(ctx, c1tChan.ChannelID, user1.KeyName(), transfer, ibc.TransferOptions{Memo: string(memo)})
		require.NoError(t, err)
		_, err = testutil.PollForAck(ctx, chain1, chain1Height, chain1Height+30, transferTx.Packet)
		require.NoError(t, err)
		err = testutil.WaitForBlocks(ctx, 1, chain1)
		require.NoError(t, err)

		chain1Balance, err := chain1.GetBalance(ctx, user1.FormattedAddress(), t1IBCDenom)
		require.NoError(t, err)
		chainTBalance, err := chainT.GetBalance(ctx, userT.FormattedAddress(), chainT.Config().Denom)
		require.NoError(t, err)
		chain2Balance, err := chain2.GetBalance(ctx, user2.FormattedAddress(), t2IBCDenom)
		require.NoError(t, err)

		require.Equal(t, chain1OldBalance, chain1Balance)
		require.Equal(t, chainTOld2Balance, chainTBalance)
		require.Equal(t, chain2OldBalance, chain2Balance)
	})
}
