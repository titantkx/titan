package app

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	snapshottypes "github.com/cosmos/cosmos-sdk/snapshots/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdksimtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/titantkx/ethermint/crypto/ethsecp256k1"
	"github.com/titantkx/titan/app/params"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/utils"
)

// SetupOptions defines arguments that are passed into `Simapp` constructor.
type SetupOptions struct {
	Logger  log.Logger
	DB      *dbm.MemDB
	AppOpts servertypes.AppOptions
}

type SnapshotsConfig struct {
	blocks             uint64
	blockTxs           int
	snapshotInterval   uint64
	snapshotKeepRecent uint32
	pruningOpts        pruningtypes.PruningOptions
}

var DefaultChainID = fmt.Sprintf("%s-1", utils.MainnetChainID)

func setup(withGenesis bool, invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp)) (*App, GenesisState, params.EncodingConfig) {
	db := dbm.NewMemDB()

	appOptions := make(sdksimtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = invCheckPeriod

	encodingConfig := MakeEncodingConfig()

	app := New(log.NewNopLogger(), db, nil, true, map[int64]bool{},
		DefaultNodeHome,
		0,
		encodingConfig,
		appOptions,
		baseAppOptions...,
	)
	if withGenesis {
		return app, NewDefaultGenesisState(app.AppCodec()), encodingConfig
	}
	return app, GenesisState{}, encodingConfig
}

// Main Setup new App
//
//

// NewSimappWithCustomOptions initializes a new SimApp with custom options.
func NewSimappWithCustomOptions(t *testing.T, isCheckTx bool, options SetupOptions) *App {
	t.Helper()

	utils.InitSDKConfig()
	utils.RegisterDenoms()

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
		Coins:   sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(1e8).Mul(sdk.NewInt(1e18)))),
	}

	app := New(options.Logger, options.DB, nil, true, map[int64]bool{},
		DefaultNodeHome,
		0,
		MakeEncodingConfig(), options.AppOpts)
	genesisState := NewDefaultGenesisState(app.AppCodec())
	genesisState, err = sdksimtestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, []authtypes.GenesisAccount{acc}, balance)
	require.NoError(t, err)

	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := tmjson.MarshalIndent(genesisState, "", " ")
		require.NoError(t, err)

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: sdksimtestutil.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp. Return app and genesis address.
func Setup(t *testing.T, isCheckTx bool) (*App, sdk.AccAddress) {
	t.Helper()

	utils.InitSDKConfig()
	utils.RegisterDenoms()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)

	genAcc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: genAcc.GetAddress().String(),
		Coins: sdk.NewCoins(
			sdk.NewCoin(utils.BaseDenom, sdk.NewInt(1e8).Mul(sdk.NewInt(1e18))),
			sdk.NewCoin(testutil.SecondaryDenom, sdk.NewInt(1e8).Mul(sdk.NewInt(1e18))),
		),
	}
	app := SetupWithGenesisValSet(t, isCheckTx, valSet, []authtypes.GenesisAccount{genAcc}, balance)

	return app, genAcc.GetAddress()
}

//nolint:unused
func setTxSignature(t *testing.T, builder client.TxBuilder, nonce uint64) {
	privKey := secp256k1.GenPrivKeyFromSecret([]byte("test"))
	pubKey := privKey.PubKey()
	err := builder.SetSignatures(
		signingtypes.SignatureV2{
			PubKey:   pubKey,
			Sequence: nonce,
			Data:     &signingtypes.SingleSignatureData{},
		},
	)
	require.NoError(t, err)
}

func SetupWithSnapshot(t *testing.T, cfg SnapshotsConfig,
	valSet *tmtypes.ValidatorSet,
	acc []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) *App {
	t.Helper()

	utils.InitSDKConfig()
	utils.RegisterDenoms()

	snapshotTimeout := 1 * time.Minute
	snapshotStore, err := snapshots.NewStore(dbm.NewMemDB(), sdktestutil.GetTempDir(t))
	require.NoError(t, err)

	app, genesisState, _ := setup(true, 5,
		baseapp.SetSnapshot(snapshotStore, snapshottypes.NewSnapshotOptions(cfg.snapshotInterval, cfg.snapshotKeepRecent)),
		baseapp.SetPruning(cfg.pruningOpts),
		baseapp.SetChainID(DefaultChainID),
	)
	genesisStateWithValSet, err := sdksimtestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, acc, balances...)
	require.NoError(t, err)

	stateBytes, err := json.MarshalIndent(genesisStateWithValSet, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			ChainId:         DefaultChainID,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: sdksimtestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	// app.Commit()

	// r := rand.New(rand.NewSource(3920758213583))
	// keyCounter := 0

	for height := uint64(1); height <= cfg.blocks; height++ {
		currentBlockHeight := app.LastBlockHeight() + 1
		app.Logger().Debug("Creating block", "height", currentBlockHeight)

		app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
			ChainID:            DefaultChainID,
			Height:             currentBlockHeight,
			AppHash:            app.LastCommitID().Hash,
			ValidatorsHash:     valSet.Hash(),
			NextValidatorsHash: valSet.Hash(),
		}})

		// for txNum := 0; txNum < cfg.blockTxs; txNum++ {
		// msgs := []sdk.Msg{}
		// for msgNum := 0; msgNum < 100; msgNum++ {
		// 	key := []byte(fmt.Sprintf("%v", keyCounter))
		// 	value := make([]byte, 10000)

		// 	_, err := r.Read(value)
		// 	require.NoError(t, err)

		// 	msgs = append(msgs, &baseapptestutil.MsgKeyValue{Key: key, Value: value})
		// 	keyCounter++
		// }

		// builder := encodingConfig.TxConfig.NewTxBuilder()
		// builder.SetMsgs(msgs...)
		// setTxSignature(t, builder, 0)

		// txBytes, err := encodingConfig.TxConfig.TxEncoder()(builder.GetTx())
		// require.NoError(t, err)

		// resp := app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
		// require.True(t, resp.IsOK(), "%v", resp.String())
		// }

		app.EndBlock(abci.RequestEndBlock{Height: currentBlockHeight})

		app.Commit()

		// wait for snapshot to be taken, since it happens asynchronously
		if cfg.snapshotInterval > 0 && height%cfg.snapshotInterval == 0 {
			start := time.Now()
			for {
				if time.Since(start) > snapshotTimeout {
					t.Errorf("timed out waiting for snapshot after %v", snapshotTimeout)
				}

				snapshot, err := snapshotStore.Get(height, snapshottypes.CurrentFormat)
				require.NoError(t, err)

				if snapshot != nil {
					break
				}

				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	return app
}

// Utility functions
//
//

// SetupWithGenesisValSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func SetupWithGenesisValSet(t *testing.T, isCheckTx bool, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *App {
	t.Helper()

	app, genesisState, _ := setup(true, 5, baseapp.SetChainID(DefaultChainID))
	genesisState, err := sdksimtestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, genAccs, balances...)
	require.NoError(t, err)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	if !isCheckTx {
		// init chain will set the validator set and initialize the genesis accounts
		app.InitChain(
			abci.RequestInitChain{
				ChainId:         DefaultChainID,
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: sdksimtestutil.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)

		// commit genesis changes
		app.Commit()
		app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
			ChainID:            DefaultChainID,
			Height:             app.LastBlockHeight() + 1,
			AppHash:            app.LastCommitID().Hash,
			ValidatorsHash:     valSet.Hash(),
			NextValidatorsHash: valSet.Hash(),
		}})
	}

	return app
}

// GenesisStateWithSingleValidator initializes GenesisState with a single validator and genesis accounts
// that also act as delegators.
func GenesisStateWithSingleValidator(t *testing.T, app *App) GenesisState {
	t.Helper()

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
	balances := []banktypes.Balance{
		{
			Address: acc.GetAddress().String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
		},
	}

	genesisState := NewDefaultGenesisState(app.AppCodec())
	genesisState, err = sdksimtestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, []authtypes.GenesisAccount{acc}, balances...)
	require.NoError(t, err)

	return genesisState
}

// NewTestNetworkFixture returns a new simapp AppConstructor for network simulation tests
func NewTestNetworkFixture() network.TestFixture {
	dir, err := os.MkdirTemp("", "simapp")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)

	app := New(log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{},
		DefaultNodeHome,
		0,
		MakeEncodingConfig(), sdksimtestutil.NewAppOptionsWithFlagHome(dir))

	appCtr := func(val network.ValidatorI) servertypes.Application {
		return New(
			val.GetCtx().Logger, dbm.NewMemDB(), nil, true,
			map[int64]bool{},
			DefaultNodeHome,
			0,
			MakeEncodingConfig(),
			sdksimtestutil.NewAppOptionsWithFlagHome(val.GetCtx().Config.RootDir),
			baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			baseapp.SetChainID(val.GetCtx().Viper.GetString(flags.FlagChainID)),
		)
	}

	return network.TestFixture{
		AppConstructor: appCtr,
		GenesisState:   NewDefaultGenesisState(app.AppCodec()),
		EncodingConfig: moduletestutil.TestEncodingConfig{
			InterfaceRegistry: app.InterfaceRegistry(),
			Codec:             app.AppCodec(),
			TxConfig:          app.TxConfig(),
			Amino:             app.LegacyAmino(),
		},
	}
}

func PrintExported(exportedApp servertypes.ExportedApp) {
	var doc tmtypes.GenesisDoc
	doc.AppState = exportedApp.AppState
	doc.Validators = exportedApp.Validators
	doc.InitialHeight = exportedApp.Height
	doc.ConsensusParams = &tmtypes.ConsensusParams{
		Block: tmtypes.BlockParams{
			MaxBytes: exportedApp.ConsensusParams.Block.MaxBytes,
			MaxGas:   exportedApp.ConsensusParams.Block.MaxGas,
		},
		Evidence: tmtypes.EvidenceParams{
			MaxAgeNumBlocks: exportedApp.ConsensusParams.Evidence.MaxAgeNumBlocks,
			MaxAgeDuration:  exportedApp.ConsensusParams.Evidence.MaxAgeDuration,
			MaxBytes:        exportedApp.ConsensusParams.Evidence.MaxBytes,
		},
		Validator: tmtypes.ValidatorParams{
			PubKeyTypes: exportedApp.ConsensusParams.Validator.PubKeyTypes,
		},
	}

	encoded, _ := tmjson.Marshal(doc)
	out := encoded

	var exportedGenDoc tmtypes.GenesisDoc
	err := tmjson.Unmarshal(out, &exportedGenDoc)
	if err != nil {
		fmt.Println("err", err)
	}
	genDocBytes, _ := tmjson.MarshalIndent(exportedGenDoc, "", "  ")
	fmt.Println("exportedGenDoc", string(genDocBytes))
}

type KeeperTestHelper struct {
	suite.Suite
}

func (s *KeeperTestHelper) SkipIfWSL() {
	SkipIfWSL(s.T())
}

// SkipIfWSL skips tests if running on WSL
// This is a workaround to enable quickly running full unit test suite locally
// on WSL without failures. The failures are stemming from trying to upload
// wasm code. An OS permissioning issue.
func SkipIfWSL(t *testing.T) {
	t.Helper()
	skip := os.Getenv("SKIP_WASM_WSL_TESTS")
	if skip == "true" {
		t.Skip("Skipping Wasm tests")
	}
}

// AssertEventEmitted asserts that ctx's event manager has emitted the given number of events
// of the given type.
func (s *KeeperTestHelper) AssertEventEmitted(ctx sdk.Context, eventTypeExpected string, numEventsExpected int) {
	allEvents := ctx.EventManager().Events()
	// filter out other events
	actualEvents := make([]sdk.Event, 0)
	for _, event := range allEvents {
		if event.Type == eventTypeExpected {
			actualEvents = append(actualEvents, event)
		}
	}
	s.Require().Equal(numEventsExpected, len(actualEvents))
}
