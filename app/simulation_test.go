package app_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"cosmossdk.io/math"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simulationtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/ethermint/x/feemarket"
	feemarkettypes "github.com/titantkx/ethermint/x/feemarket/types"
	"github.com/titantkx/titan/app"
	nftminttypes "github.com/titantkx/titan/x/nftmint/types"
	validatorrewardtypes "github.com/titantkx/titan/x/validatorreward/types"
)

var ModuleBasics = app.ModuleBasics

// Override feemarket.AppModuleBasic
type FeeMarketAppModuleBasic struct {
	feemarket.AppModuleBasic
}

func (FeeMarketAppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	params := feemarkettypes.DefaultParams()
	params.NoBaseFee = true
	params.BaseFee = math.NewInt(0)
	params.MinGasPrice = math.LegacyNewDec(0)

	genState := &feemarkettypes.GenesisState{
		Params:   params,
		BlockGas: 0,
	}

	return cdc.MustMarshalJSON(genState)
}

type storeKeysPrefixes struct {
	A        storetypes.StoreKey
	B        storetypes.StoreKey
	Prefixes [][]byte
}

// Get flags every time the simulator is run
func init() {
	simcli.GetSimulatorFlags()

	ModuleBasics[FeeMarketAppModuleBasic{}.Name()] = FeeMarketAppModuleBasic{}
}

// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// an IAVLStore for faster simulation speed.
func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

// BenchmarkSimulation run the chain simulation
// Running using starport command:
// `starport chain simulate -v --numBlocks 200 --blockSize 50`
// Running as go benchmark test:
// `go test -benchmem -run=^$ -bench=^BenchmarkSimulation$ ./app -NumBlocks=200 -BlockSize 50 -Commit=true -Verbose=true -Enabled=true`
func BenchmarkSimulation(b *testing.B) {
	simcli.FlagSeedValue = time.Now().Unix()
	simcli.FlagVerboseValue = true
	simcli.FlagCommitValue = true
	simcli.FlagEnabledValue = true

	config := simcli.NewConfigFromFlags()
	config.ChainID = "titanben_1-1"
	db, dir, logger, _, err := simtestutil.SetupSimulation(
		config,
		"leveldb-bApp-sim",
		"Simulation",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)
	require.NoError(b, err, "simulation setup failed")

	b.Cleanup(func() {
		require.NoError(b, db.Close())
		require.NoError(b, os.RemoveAll(dir))
	})

	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = app.DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	bApp := app.New(
		logger,
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		app.MakeEncodingConfig(),
		appOptions,
		baseapp.SetChainID(config.ChainID),
	)
	require.Equal(b, app.Name, bApp.Name())

	vestingypes.RegisterLegacyAminoCodec(bApp.LegacyAmino())
	vestingypes.RegisterInterfaces(bApp.InterfaceRegistry())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		b,
		os.Stdout,
		bApp.BaseApp,
		simtestutil.AppStateFn(
			bApp.AppCodec(),
			bApp.SimulationManager(),
			ModuleBasics.DefaultGenesis(bApp.AppCodec()),
		),
		simulationtypes.RandomAccounts,
		simtestutil.SimulationOperations(bApp, bApp.AppCodec(), config),
		bApp.ModuleAccountAddrs(),
		config,
		bApp.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simtestutil.CheckExportSimulation(bApp, config, simParams)
	require.NoError(b, err)
	require.NoError(b, simErr)

	if config.Commit {
		simtestutil.PrintStats(db)
	}
}

// `go test -benchmem -run=^TestAppStateDeterminism$ -bench ^$ ./app -NumBlocks=200 -BlockSize 50 -Commit=true -Verbose=true -Enabled=true`
func TestAppStateDeterminism(t *testing.T) {
	if !simcli.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simcli.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = true
	config.AllInvariants = true

	var (
		r                    = rand.New(rand.NewSource(time.Now().Unix()))
		numSeeds             = 3
		numTimesToRunPerSeed = 5
		appHashList          = make([]json.RawMessage, numTimesToRunPerSeed)
		appOptions           = make(simtestutil.AppOptionsMap, 0)
	)
	appOptions[flags.FlagHome] = app.DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	for i := 0; i < numSeeds; i++ {
		config.Seed = r.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if simcli.FlagVerboseValue {
				logger = log.TestingLogger()
			} else {
				logger = log.NewNopLogger()
			}
			chainID := fmt.Sprintf("titan_%d-%d", i+1, j+1)
			config.ChainID = chainID

			db := dbm.NewMemDB()
			bApp := app.New(
				logger,
				db,
				nil,
				true,
				map[int64]bool{},
				app.DefaultNodeHome,
				simcli.FlagPeriodValue,
				app.MakeEncodingConfig(),
				appOptions,
				fauxMerkleModeOpt,
				baseapp.SetChainID(chainID),
			)

			vestingypes.RegisterLegacyAminoCodec(bApp.LegacyAmino())
			vestingypes.RegisterInterfaces(bApp.InterfaceRegistry())

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t,
				os.Stdout,
				bApp.BaseApp,
				simtestutil.AppStateFn(
					bApp.AppCodec(),
					bApp.SimulationManager(),
					ModuleBasics.DefaultGenesis(bApp.AppCodec()),
				),
				simulationtypes.RandomAccounts,
				simtestutil.SimulationOperations(bApp, bApp.AppCodec(), config),
				bApp.ModuleAccountAddrs(),
				config,
				bApp.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				simtestutil.PrintStats(db)
			}

			appHash := bApp.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, string(appHashList[0]), string(appHashList[j]),
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}

// `go test -benchmem -run=^TestAppImportExport$ -bench ^$ ./app -NumBlocks=200 -BlockSize 50 -Commit=true -Verbose=true -Enabled=true`
func TestAppImportExport(t *testing.T) {
	config := simcli.NewConfigFromFlags()
	config.ChainID = "titanimport_1-1"

	db, dir, logger, skip, err := simtestutil.SetupSimulation(
		config,
		"leveldb-app-sim",
		"Simulation",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)
	if skip {
		t.Skip("skipping application import/export simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		require.NoError(t, db.Close())
		require.NoError(t, os.RemoveAll(dir))
	}()

	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = app.DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	bApp := app.New(
		logger,
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		app.MakeEncodingConfig(),
		appOptions,
		baseapp.SetChainID(config.ChainID),
	)
	require.Equal(t, app.Name, bApp.Name())

	vestingypes.RegisterLegacyAminoCodec(bApp.LegacyAmino())
	vestingypes.RegisterInterfaces(bApp.InterfaceRegistry())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		bApp.BaseApp,
		simtestutil.AppStateFn(
			bApp.AppCodec(),
			bApp.SimulationManager(),
			ModuleBasics.DefaultGenesis(bApp.AppCodec()),
		),
		simulationtypes.RandomAccounts,
		simtestutil.SimulationOperations(bApp, bApp.AppCodec(), config),
		bApp.BlockedModuleAccountAddrs(),
		config,
		bApp.AppCodec(),
	)
	require.NoError(t, simErr)

	// export state and simParams before the simulation error is checked
	err = simtestutil.CheckExportSimulation(bApp, config, simParams)
	require.NoError(t, err)

	if config.Commit {
		simtestutil.PrintStats(db)
	}

	fmt.Printf("exporting genesis...\n")

	exported, err := bApp.ExportAppStateAndValidators(false, []string{}, []string{})
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	newDB, newDir, _, _, err := simtestutil.SetupSimulation(
		config,
		"leveldb-app-sim-2",
		"Simulation-2",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		require.NoError(t, newDB.Close())
		require.NoError(t, os.RemoveAll(newDir))
	}()

	newApp := app.New(
		log.NewNopLogger(),
		newDB,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		app.MakeEncodingConfig(),
		appOptions,
		baseapp.SetChainID(config.ChainID),
	)
	require.Equal(t, app.Name, newApp.Name())

	vestingypes.RegisterLegacyAminoCodec(newApp.LegacyAmino())
	vestingypes.RegisterInterfaces(newApp.InterfaceRegistry())

	var genesisState app.GenesisState
	err = json.Unmarshal(exported.AppState, &genesisState)
	require.NoError(t, err)

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%v", r)
			if !strings.Contains(err, "validator set is empty after InitGenesis") {
				panic(r)
			}
			logger.Info("Skipping simulation as all validators have been unbonded")
			logger.Info("err", err, "stacktrace", string(debug.Stack()))
		}
	}()

	ctxA := bApp.NewContext(true, tmproto.Header{ChainID: config.ChainID, Height: bApp.LastBlockHeight()})
	ctxB := newApp.NewContext(true, tmproto.Header{ChainID: config.ChainID, Height: bApp.LastBlockHeight()})

	newApp.ModuleManager().InitGenesis(ctxB, bApp.AppCodec(), genesisState)
	newApp.StoreConsensusParams(ctxB, exported.ConsensusParams)

	fmt.Printf("comparing stores...\n")

	storeKeysPrefixes := []storeKeysPrefixes{
		{bApp.GetKey(authtypes.StoreKey), newApp.GetKey(authtypes.StoreKey), [][]byte{}},
		{
			bApp.GetKey(stakingtypes.StoreKey), newApp.GetKey(stakingtypes.StoreKey),
			[][]byte{
				stakingtypes.UnbondingQueueKey, stakingtypes.RedelegationQueueKey, stakingtypes.ValidatorQueueKey,
				stakingtypes.HistoricalInfoKey, stakingtypes.UnbondingIDKey, stakingtypes.UnbondingIndexKey, stakingtypes.UnbondingTypeKey, stakingtypes.ValidatorUpdatesKey,
			},
		}, // ordering may change but it doesn't matter
		{bApp.GetKey(slashingtypes.StoreKey), newApp.GetKey(slashingtypes.StoreKey), [][]byte{}},
		{bApp.GetKey(distrtypes.StoreKey), newApp.GetKey(distrtypes.StoreKey), [][]byte{}},
		{bApp.GetKey(banktypes.StoreKey), newApp.GetKey(banktypes.StoreKey), [][]byte{banktypes.BalancesPrefix}},
		{bApp.GetKey(paramstypes.StoreKey), newApp.GetKey(paramstypes.StoreKey), [][]byte{}},
		{bApp.GetKey(govtypes.StoreKey), newApp.GetKey(govtypes.StoreKey), [][]byte{}},
		{bApp.GetKey(evidencetypes.StoreKey), newApp.GetKey(evidencetypes.StoreKey), [][]byte{}},
		{bApp.GetKey(capabilitytypes.StoreKey), newApp.GetKey(capabilitytypes.StoreKey), [][]byte{}},
		{bApp.GetKey(authzkeeper.StoreKey), newApp.GetKey(authzkeeper.StoreKey), [][]byte{authzkeeper.GrantKey, authzkeeper.GrantQueuePrefix}},
		{bApp.GetKey(validatorrewardtypes.StoreKey), newApp.GetKey(validatorrewardtypes.StoreKey), [][]byte{}},
		{bApp.GetKey(nftminttypes.StoreKey), newApp.GetKey(nftminttypes.StoreKey), [][]byte{}},
	}

	for _, skp := range storeKeysPrefixes {
		storeA := ctxA.KVStore(skp.A)
		storeB := ctxB.KVStore(skp.B)

		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, skp.Prefixes)
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

		fmt.Printf("compared %d different key/value pairs between %s and %s\n", len(failedKVAs), skp.A, skp.B)
		require.Equal(t, 0, len(failedKVAs), simtestutil.GetSimulationLog(skp.A.Name(), bApp.SimulationManager().StoreDecoders, failedKVAs, failedKVBs))
	}
}

// `go test -benchmem -run=^TestAppSimulationAfterImport$ -bench ^$ ./app -NumBlocks=200 -BlockSize 50 -Commit=true -Verbose=true -Enabled=true`
func TestAppSimulationAfterImport(t *testing.T) {
	config := simcli.NewConfigFromFlags()
	config.ChainID = "titanafterimport_1-1"

	db, dir, logger, skip, err := simtestutil.SetupSimulation(
		config,
		"leveldb-app-sim",
		"Simulation",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)
	if skip {
		t.Skip("skipping application simulation after import")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		require.NoError(t, db.Close())
		require.NoError(t, os.RemoveAll(dir))
	}()

	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = app.DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	bApp := app.New(
		logger,
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		app.MakeEncodingConfig(),
		appOptions,
		fauxMerkleModeOpt,
		baseapp.SetChainID(config.ChainID),
	)
	require.Equal(t, app.Name, bApp.Name())

	vestingypes.RegisterLegacyAminoCodec(bApp.LegacyAmino())
	vestingypes.RegisterInterfaces(bApp.InterfaceRegistry())

	// run randomized simulation
	stopEarly, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		bApp.BaseApp,
		simtestutil.AppStateFn(
			bApp.AppCodec(),
			bApp.SimulationManager(),
			ModuleBasics.DefaultGenesis(bApp.AppCodec()),
		),
		simulationtypes.RandomAccounts,
		simtestutil.SimulationOperations(bApp, bApp.AppCodec(), config),
		bApp.BlockedModuleAccountAddrs(),
		config,
		bApp.AppCodec(),
	)
	require.NoError(t, simErr)

	// export state and simParams before the simulation error is checked
	err = simtestutil.CheckExportSimulation(bApp, config, simParams)
	require.NoError(t, err)

	if config.Commit {
		simtestutil.PrintStats(db)
	}

	if stopEarly {
		fmt.Println("can't export or import a zero-validator genesis, exiting test...")
		return
	}

	fmt.Printf("exporting genesis...\n")

	exported, err := bApp.ExportAppStateAndValidators(true, []string{}, []string{})
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	newDB, newDir, _, _, err := simtestutil.SetupSimulation(
		config,
		"leveldb-app-sim-2",
		"Simulation-2",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		require.NoError(t, newDB.Close())
		require.NoError(t, os.RemoveAll(newDir))
	}()

	newApp := app.New(
		log.NewNopLogger(),
		newDB,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		app.MakeEncodingConfig(),
		appOptions,
		fauxMerkleModeOpt,
		baseapp.SetChainID(config.ChainID),
	)
	require.Equal(t, app.Name, newApp.Name())

	vestingypes.RegisterLegacyAminoCodec(newApp.LegacyAmino())
	vestingypes.RegisterInterfaces(newApp.InterfaceRegistry())

	newApp.InitChain(abci.RequestInitChain{
		ChainId:       config.ChainID,
		AppStateBytes: exported.AppState,
	})

	_, _, err = simulation.SimulateFromSeed(
		t,
		os.Stdout,
		newApp.BaseApp,
		simtestutil.AppStateFn(
			bApp.AppCodec(),
			bApp.SimulationManager(),
			ModuleBasics.DefaultGenesis(newApp.AppCodec()),
		),
		simulationtypes.RandomAccounts,
		simtestutil.SimulationOperations(newApp, newApp.AppCodec(), config),
		newApp.BlockedModuleAccountAddrs(),
		config,
		bApp.AppCodec(),
	)
	require.NoError(t, err)
}
