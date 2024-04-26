package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkstakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/titantkx/titan/app"
)

type IntegrationTestSuite struct {
	suite.Suite

	app *app.App
	ctx sdk.Context

	genAddr sdk.AccAddress
	addrs   []sdk.AccAddress
	vals    []types.Validator

	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func (s *IntegrationTestSuite) SetupTest() {
	s.app, s.genAddr = app.Setup(s.T(), false)
	ctx := s.app.BaseApp.NewContext(false, tmproto.Header{})

	querier := sdkstakingkeeper.Querier{Keeper: s.app.StakingKeeper.Keeper}

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, s.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)
	queryClient := types.NewQueryClient(queryHelper)

	s.msgServer = sdkstakingkeeper.NewMsgServerImpl(s.app.StakingKeeper.Keeper)

	addrs, _, validators := createValidators(s.T(), ctx, s.app, s.genAddr, []int64{9, 8, 7})

	header := tmproto.Header{
		ChainID: "HelloChain",
		Height:  5,
	}

	// sort a copy of the validators, so that original validators does not
	// have its order changed
	sortedVals := make([]types.Validator, len(validators))
	copy(sortedVals, validators)
	hi := types.NewHistoricalInfo(header, sortedVals, s.app.StakingKeeper.PowerReduction(ctx))
	s.app.StakingKeeper.SetHistoricalInfo(ctx, 5, &hi)

	s.ctx, s.queryClient = ctx, queryClient
	s.addrs, s.vals = addrs, validators
}

// NewValidator is a testing helper method to create validators in tests
func NewValidator(suite *IntegrationTestSuite, operator sdk.ValAddress, pubKey cryptotypes.PubKey) types.Validator {
	v, err := types.NewValidator(operator, pubKey, types.Description{})
	suite.NoError(err)
	return v
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
