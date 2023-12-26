package keeper_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/golang/mock/gomock"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	sdkgovkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	sdkgovtestutil "github.com/cosmos/cosmos-sdk/x/gov/testutil"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	v1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/tokenize-titan/titan/utils"
	"github.com/tokenize-titan/titan/x/gov/keeper"
	govtestutil "github.com/tokenize-titan/titan/x/gov/testutil"
)

var (
	_, _, addr                 = testdata.KeyTestPubAddr()
	_, _, genAddr              = testdata.KeyTestPubAddr()
	govAcct                    = sdk.AccAddress{}
	TestProposal               = []sdk.Msg{}
	MinDepositByConcensusPower = int64(100)
)

func initGlobalVars() {
	govAcct = authtypes.NewModuleAddress(govtypes.ModuleName)
	TestProposal = getTestProposal()
}

// getTestProposal creates and returns a test proposal message.
func getTestProposal() []sdk.Msg {
	legacyProposalMsg, err := v1.NewLegacyContent(v1beta1.NewTextProposal("Title", "description"), authtypes.NewModuleAddress(govtypes.ModuleName).String())
	if err != nil {
		panic(err)
	}

	return []sdk.Msg{
		banktypes.NewMsgSend(govAcct, addr, sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(250).Mul(sdk.NewInt(1e18))))),
		legacyProposalMsg,
	}
}

func setupGovKeeper(t *testing.T) (
	*keeper.Keeper,
	*sdkgovtestutil.MockAccountKeeper,
	*sdkgovtestutil.MockBankKeeper,
	*sdkgovtestutil.MockStakingKeeper,
	*govtestutil.MockDistributionKeeper,
	sdk.Address,
	moduletestutil.TestEncodingConfig,
	sdk.Context,
) {
	t.Helper()

	// fmt.Println(addr.String(), genAddr.String(), govAcct.String())
	utils.InitSDKConfig()
	utils.RegisterDenoms()

	initGlobalVars()

	govAcct = authtypes.NewModuleAddress(govtypes.ModuleName)

	key := sdk.NewKVStoreKey(govtypes.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, sdk.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()
	v1.RegisterInterfaces(encCfg.InterfaceRegistry)
	v1beta1.RegisterInterfaces(encCfg.InterfaceRegistry)
	banktypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	// Create MsgServiceRouter, but don't populate it before creating the gov
	// keeper.
	msr := baseapp.NewMsgServiceRouter()

	// init mock
	ctrl := gomock.NewController(t)
	accKeeper := sdkgovtestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := sdkgovtestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := sdkgovtestutil.NewMockStakingKeeper(ctrl)
	distrKeeper := govtestutil.NewMockDistributionKeeper(ctrl)

	trackMockAccount(ctx, accKeeper)
	trackMockBank(ctx, bankKeeper)
	trackMockStaking(ctx, stakingKeeper)
	trackMockDistribution(ctx, distrKeeper, bankKeeper)

	// init gov
	govKeeper := keeper.NewKeeper(encCfg.Codec, key, accKeeper, bankKeeper, stakingKeeper, distrKeeper, msr, govtypes.DefaultConfig(), govAcct.String())
	govKeeper.SetProposalID(ctx, 1)
	govRouter := v1beta1.NewRouter() // Also register legacy gov handlers to test them too.
	govRouter.AddRoute(govtypes.RouterKey, v1beta1.ProposalHandler)
	govKeeper.SetLegacyRouter(govRouter)
	govParams := v1.DefaultParams()
	govParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(MinDepositByConcensusPower).Mul(sdk.NewInt(1e18))))
	govKeeper.SetParams(ctx, govParams)

	// Register all handlers for the MegServiceRouter.
	msr.SetInterfaceRegistry(encCfg.InterfaceRegistry)
	v1.RegisterMsgServer(msr, sdkgovkeeper.NewMsgServerImpl(govKeeper.Keeper))
	banktypes.RegisterMsgServer(msr, nil) // Nil is fine here as long as we never execute the proposal's Msgs.

	return govKeeper, accKeeper, bankKeeper, stakingKeeper, distrKeeper, genAddr, encCfg, ctx
}

func trackMockAccount(ctx sdk.Context, accKeeper *sdkgovtestutil.MockAccountKeeper) {
	accKeeper.EXPECT().GetModuleAddress(govtypes.ModuleName).Return(govAcct).AnyTimes()
	accKeeper.EXPECT().GetModuleAccount(gomock.Any(), govtypes.ModuleName).Return(authtypes.NewEmptyModuleAccount(govtypes.ModuleName)).AnyTimes()
}

// trackMockBalances sets up expected calls on the Mock BankKeeper, and also
// locally tracks accounts balances (not modules balances).
func trackMockBank(ctx sdk.Context, bankKeeper *sdkgovtestutil.MockBankKeeper) {
	balances := make(map[string]sdk.Coins)
	balances[genAddr.String()] = sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(1e8).Mul(sdk.NewInt(1e18))))

	// We don't track module account balances.
	bankKeeper.EXPECT().BurnCoins(gomock.Any(), govtypes.ModuleName, gomock.Any()).Times(0)

	bankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ sdk.Context, senderModule, recipientModule string, coins sdk.Coins) error {
		newBalance, negative := balances[senderModule].SafeSub(coins...)
		if negative {
			return fmt.Errorf("not enough balance")
		}

		balances[senderModule] = newBalance
		balances[recipientModule] = balances[recipientModule].Add(coins...)

		return nil
	}).AnyTimes()

	// But we do track normal account balances.
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), govtypes.ModuleName, gomock.Any()).DoAndReturn(func(_ sdk.Context, sender sdk.AccAddress, recipientModule string, coins sdk.Coins) error {
		newBalance, negative := balances[sender.String()].SafeSub(coins...)
		if negative {
			return fmt.Errorf("not enough balance")
		}
		balances[sender.String()] = newBalance
		balances[recipientModule] = balances[recipientModule].Add(coins...)
		return nil
	}).AnyTimes()

	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ sdk.Context, module string, rcpt sdk.AccAddress, coins sdk.Coins) error {
		newBalance, negative := balances[module].SafeSub(coins...)
		if negative {
			return fmt.Errorf("not enough balance")
		}
		balances[module] = newBalance
		balances[rcpt.String()] = balances[rcpt.String()].Add(coins...)

		return nil
	}).AnyTimes()

	bankKeeper.EXPECT().GetAllBalances(gomock.Any(), gomock.Any()).DoAndReturn(func(_ sdk.Context, addr sdk.AccAddress) sdk.Coins {
		return balances[addr.String()]
	}).AnyTimes()

	bankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ sdk.Context, senderAddr, recipientAddr sdk.AccAddress, coins sdk.Coins) error {
		newBalance, negative := balances[senderAddr.String()].SafeSub(coins...)
		if negative {
			return fmt.Errorf("not enough balance")
		}
		balances[senderAddr.String()] = newBalance
		balances[recipientAddr.String()] = balances[recipientAddr.String()].Add(coins...)
		return nil
	}).AnyTimes()
}

func trackMockStaking(ctx sdk.Context, stakingKeeper *sdkgovtestutil.MockStakingKeeper) {
	stakingKeeper.EXPECT().TokensFromConsensusPower(ctx, gomock.Any()).DoAndReturn(func(ctx sdk.Context, power int64) math.Int {
		return sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction)
	}).AnyTimes()
	stakingKeeper.EXPECT().BondDenom(ctx).Return(utils.BaseDenom).AnyTimes()
	stakingKeeper.EXPECT().IterateBondedValidatorsByPower(gomock.Any(), gomock.Any()).AnyTimes()
	stakingKeeper.EXPECT().IterateDelegations(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	stakingKeeper.EXPECT().TotalBondedTokens(gomock.Any()).Return(math.NewInt(1e3).Mul(math.NewInt(1e18))).AnyTimes()
}

func trackMockDistribution(ctx sdk.Context, distrKeeper *govtestutil.MockDistributionKeeper, bankKeeper *sdkgovtestutil.MockBankKeeper) {
	distrKeeper.EXPECT().FundCommunityPoolFromModule(gomock.Any(), gomock.Any(), govtypes.ModuleName).DoAndReturn(func(_ sdk.Context, amount sdk.Coins, _ string) error {
		if err := bankKeeper.SendCoinsFromModuleToModule(ctx, govtypes.ModuleName, distrtypes.ModuleName, amount); err != nil {
			return err
		}
		return nil
	}).AnyTimes()
}
