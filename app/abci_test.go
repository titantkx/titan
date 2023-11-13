package app

import (
	"fmt"
	"reflect"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"

	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tokenize-titan/ethermint/crypto/ethsecp256k1"
)

func TestABCI_ApplySnapshotChunk(t *testing.T) {
	//
	//
	//

	InitSDKConfig()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)

	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
	}

	//
	//
	//
	srcCfg := SnapshotsConfig{
		blocks:             4,
		blockTxs:           10,
		snapshotInterval:   2,
		snapshotKeepRecent: 2,
		pruningOpts:        pruningtypes.NewPruningOptions(pruningtypes.PruningNothing),
	}
	srcApp := SetupWithSnapshot(t, srcCfg, valSet, []authtypes.GenesisAccount{acc}, balance)

	exportedSrc, err := srcApp.ExportAppStateAndValidators(true, []string{}, []string{})
	require.NoError(t, err)

	targetCfg := SnapshotsConfig{
		blocks:             0,
		blockTxs:           10,
		snapshotInterval:   2,
		snapshotKeepRecent: 2,
		pruningOpts:        pruningtypes.NewPruningOptions(pruningtypes.PruningNothing),
	}
	targetApp := SetupWithSnapshot(t, targetCfg, valSet, []authtypes.GenesisAccount{acc}, balance)

	// fetch latest snapshot to restore
	respList := srcApp.ListSnapshots(abci.RequestListSnapshots{})
	require.NotEmpty(t, respList.Snapshots)
	snapshot := respList.Snapshots[0]
	fmt.Println("snapshot", snapshot)

	// make sure the snapshot has at least 1 chunks
	require.GreaterOrEqual(t, snapshot.Chunks, uint32(1), "Not enough snapshot chunks")

	// begin a snapshot restoration in the target
	respOffer := targetApp.OfferSnapshot(abci.RequestOfferSnapshot{Snapshot: snapshot})

	fmt.Println("respOffer", respOffer)

	require.Equal(t, abci.ResponseOfferSnapshot{Result: abci.ResponseOfferSnapshot_ACCEPT}, respOffer)

	// We should be able to pass an invalid chunk and get a verify failure, before
	// reapplying it.
	respApply := targetApp.ApplySnapshotChunk(abci.RequestApplySnapshotChunk{
		Index:  0,
		Chunk:  []byte{9},
		Sender: "sender",
	})
	require.Equal(t, abci.ResponseApplySnapshotChunk{
		Result:        abci.ResponseApplySnapshotChunk_RETRY,
		RefetchChunks: []uint32{0},
		RejectSenders: []string{"sender"},
	}, respApply)

	// fetch each chunk from the source and apply it to the target
	for index := uint32(0); index < snapshot.Chunks; index++ {
		respChunk := srcApp.LoadSnapshotChunk(abci.RequestLoadSnapshotChunk{
			Height: snapshot.Height,
			Format: snapshot.Format,
			Chunk:  index,
		})
		require.NotNil(t, respChunk.Chunk)

		respApply := targetApp.ApplySnapshotChunk(abci.RequestApplySnapshotChunk{
			Index: index,
			Chunk: respChunk.Chunk,
		})
		require.Equal(t, abci.ResponseApplySnapshotChunk{
			Result: abci.ResponseApplySnapshotChunk_ACCEPT,
		}, respApply)
	}

	exportedTarget, errExportTarget := targetApp.ExportAppStateAndValidators(true, []string{}, []string{})
	require.NoError(t, errExportTarget)

	require.True(t, reflect.DeepEqual(exportedSrc.AppState, exportedTarget.AppState))

	// the target should now have the same hash as the source
	require.Equal(t, srcApp.LastCommitID(), targetApp.LastCommitID())
}
