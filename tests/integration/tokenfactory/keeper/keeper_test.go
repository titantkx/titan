package keeper_test

import (
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/titantkx/titan/app"
	"github.com/titantkx/titan/x/tokenfactory/keeper"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

type KeeperTestSuite struct {
	app.KeeperTestHelper

	app *app.App
	ctx sdk.Context

	genAddr sdk.AccAddress

	queryClient    types.QueryClient
	msgServer      types.MsgServer
	contractKeeper wasmtypes.ContractOpsKeeper
	bankMsgServer  banktypes.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app, s.genAddr = app.Setup(s.T(), false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, s.app.TokenfactoryKeeper)

	s.contractKeeper = wasmkeeper.NewGovPermissionKeeper(s.app.WasmKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = keeper.NewMsgServerImpl(s.app.TokenfactoryKeeper)
	s.bankMsgServer = bankkeeper.NewMsgServerImpl(s.app.BankKeeper)
}

func (s *KeeperTestSuite) SetupTestForInitGenesis() {
	// Setting to True, leads to init genesis not running
	s.app, s.genAddr = app.Setup(s.T(), true)
	s.ctx = s.app.BaseApp.NewContext(true, tmproto.Header{Time: time.Now()})
}

func (s *KeeperTestSuite) CreateDenom(creator string, denom string) string {
	res, err := s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(creator, denom))
	s.Require().NoError(err)

	return res.GetNewTokenDenom()
}

func (s *KeeperTestSuite) FundAcc(addr sdk.AccAddress, coins sdk.Coins) {
	err := s.app.BankKeeper.SendCoins(s.ctx, s.genAddr, addr, coins)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCreateModuleAccount() {
	app := s.app

	// setup new next account number
	nextAccountNumber := app.AccountKeeper.NextAccountNumber(s.ctx)

	// remove module account
	tokenfactoryModuleAccount := app.AccountKeeper.GetAccount(s.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	app.AccountKeeper.RemoveAccount(s.ctx, tokenfactoryModuleAccount)

	// ensure module account was removed
	s.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	tokenfactoryModuleAccount = app.AccountKeeper.GetAccount(s.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().Nil(tokenfactoryModuleAccount)

	// create module account
	app.TokenfactoryKeeper.CreateModuleAccount(s.ctx)

	// check that the module account is now initialized
	tokenfactoryModuleAccount = app.AccountKeeper.GetAccount(s.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().NotNil(tokenfactoryModuleAccount)

	// check that the account number of the module account is now initialized correctly
	s.Require().Equal(nextAccountNumber+1, tokenfactoryModuleAccount.GetAccountNumber())
}
