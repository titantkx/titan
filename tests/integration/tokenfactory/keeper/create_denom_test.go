package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func (s *KeeperTestSuite) TestMsgCreateDenom() {
	var (
		tokenFactoryKeeper = s.app.TokenfactoryKeeper
		bankKeeper         = s.app.BankKeeper
		denomCreationFee   = sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, math.NewInt(1000000)))
	)

	// Set the denom creation fee. It is currently turned off in favor
	// of gas charge by default.
	params := s.app.TokenfactoryKeeper.GetParams(s.ctx)
	params.DenomCreationFee = denomCreationFee
	err := s.app.TokenfactoryKeeper.SetParams(s.ctx, params)
	s.Require().NoError(err)

	addr0 := sample.AccAddress()
	addr1 := sample.AccAddress()

	// Fund denom creation fee for every execution of MsgCreateDenom.
	s.FundAcc(addr0, denomCreationFee)
	s.FundAcc(addr0, denomCreationFee)
	s.FundAcc(addr1, denomCreationFee)

	// Get balance of acc 0 before creating a denom
	preCreateBalance := bankKeeper.GetBalance(s.ctx, addr0, denomCreationFee[0].Denom)

	// Creating a denom should work
	res, err := s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(addr0.String(), "bitcoin"))
	s.Require().NoError(err)
	s.Require().NotEmpty(res.GetNewTokenDenom())

	// Make sure that the admin is set correctly
	queryRes, err := s.queryClient.DenomAuthorityMetadata(s.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
		Denom: res.GetNewTokenDenom(),
	})
	s.Require().NoError(err)
	s.Require().Equal(addr0.String(), queryRes.AuthorityMetadata.Admin)

	// Make sure that creation fee was deducted
	postCreateBalance := bankKeeper.GetBalance(s.ctx, addr0, tokenFactoryKeeper.GetParams(s.ctx).DenomCreationFee[0].Denom)
	s.Require().True(preCreateBalance.Sub(postCreateBalance).IsEqual(denomCreationFee[0]))

	// Make sure that a second version of the same denom can't be recreated
	_, err = s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(addr0.String(), "bitcoin"))
	s.Require().Error(err)

	// Creating a second denom should work
	res, err = s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(addr0.String(), "litecoin"))
	s.Require().NoError(err)
	s.Require().NotEmpty(res.GetNewTokenDenom())

	// Try querying all the denoms created by addr0
	queryRes2, err := s.queryClient.DenomsFromCreator(s.ctx.Context(), &types.QueryDenomsFromCreatorRequest{
		Creator: addr0.String(),
	})
	s.Require().NoError(err)
	s.Require().Len(queryRes2.Denoms, 2)

	// Make sure that a second account can create a denom with the same subdenom
	res, err = s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(addr1.String(), "bitcoin"))
	s.Require().NoError(err)
	s.Require().NotEmpty(res.GetNewTokenDenom())

	// Make sure that an address with a "/" in it can't create denoms
	_, err = s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom("titan.eth/creator", "bitcoin"))
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestCreateDenom() {
	var (
		primaryDenom            = utils.BaseDenom
		secondaryDenom          = testutil.SecondaryDenom
		defaultDenomCreationFee = types.Params{DenomCreationFee: sdk.NewCoins(sdk.NewCoin(primaryDenom, math.NewInt(50000000)))}
		twoDenomCreationFee     = types.Params{DenomCreationFee: sdk.NewCoins(sdk.NewCoin(primaryDenom, math.NewInt(50000000)), sdk.NewCoin(secondaryDenom, math.NewInt(50000000)))}
		nilCreationFee          = types.Params{DenomCreationFee: nil}
		largeCreationFee        = types.Params{DenomCreationFee: sdk.NewCoins(sdk.NewCoin(primaryDenom, math.NewInt(5000000000)))}
	)

	addr0 := sample.AccAddress()

	for _, tc := range []struct {
		desc             string
		fund             sdk.Coins
		denomCreationFee types.Params
		setup            func()
		subdenom         string
		valid            bool
	}{
		{
			desc:             "subdenom too long",
			fund:             defaultDenomCreationFee.DenomCreationFee,
			denomCreationFee: defaultDenomCreationFee,
			subdenom:         "assadsadsadasdasdsadsadsadsadsadsadsklkadaskkkdasdasedskhanhassyeunganassfnlksdflksafjlkasd",
			valid:            false,
		},
		{
			desc:             "subdenom and creator pair already exists",
			fund:             defaultDenomCreationFee.DenomCreationFee,
			denomCreationFee: defaultDenomCreationFee,
			setup: func() {
				_, err := s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(addr0.String(), "bitcoin"))
				s.Require().NoError(err)
			},
			subdenom: "bitcoin",
			valid:    false,
		},
		{
			desc:             "success case: defaultDenomCreationFee",
			fund:             defaultDenomCreationFee.DenomCreationFee,
			denomCreationFee: defaultDenomCreationFee,
			subdenom:         "evmos",
			valid:            true,
		},
		{
			desc:             "success case: twoDenomCreationFee",
			fund:             twoDenomCreationFee.DenomCreationFee,
			denomCreationFee: twoDenomCreationFee,
			subdenom:         "catcoin",
			valid:            true,
		},
		{
			desc:             "success case: nilCreationFee",
			fund:             sdk.NewCoins(),
			denomCreationFee: nilCreationFee,
			subdenom:         "czcoin",
			valid:            true,
		},
		{
			desc:             "account doesn't have enough to pay for denom creation fee",
			fund:             defaultDenomCreationFee.DenomCreationFee,
			denomCreationFee: largeCreationFee,
			subdenom:         "tooexpensive",
			valid:            false,
		},
		{
			desc:             "subdenom having invalid characters",
			fund:             defaultDenomCreationFee.DenomCreationFee,
			denomCreationFee: defaultDenomCreationFee,
			subdenom:         "bit/***///&&&/coin",
			valid:            false,
		},
	} {
		s.SetupTest()
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			if tc.setup != nil {
				tc.setup()
			}
			tokenFactoryKeeper := s.app.TokenfactoryKeeper
			bankKeeper := s.app.BankKeeper
			// Set denom creation fee in params
			err := tokenFactoryKeeper.SetParams(s.ctx, tc.denomCreationFee)
			s.Require().NoError(err)
			denomCreationFee := tokenFactoryKeeper.GetParams(s.ctx).DenomCreationFee
			s.Require().Equal(tc.denomCreationFee.DenomCreationFee, denomCreationFee)

			s.FundAcc(addr0, tc.fund)

			// note balance, create a tokenfactory denom, then note balance again
			preCreateBalance := bankKeeper.GetAllBalances(s.ctx, addr0)
			res, err := s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(addr0.String(), tc.subdenom))
			postCreateBalance := bankKeeper.GetAllBalances(s.ctx, addr0)
			if tc.valid {
				s.Require().NoError(err)
				s.Require().True(preCreateBalance.Sub(postCreateBalance...).IsEqual(denomCreationFee))

				// Make sure that the admin is set correctly
				queryRes, err := s.queryClient.DenomAuthorityMetadata(s.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
					Denom: res.GetNewTokenDenom(),
				})

				s.Require().NoError(err)
				s.Require().Equal(addr0.String(), queryRes.AuthorityMetadata.Admin)

				// Make sure that the denom metadata is initialized correctly
				metadata, found := bankKeeper.GetDenomMetaData(s.ctx, res.GetNewTokenDenom())
				s.Require().True(found)
				s.Require().Equal(banktypes.Metadata{
					DenomUnits: []*banktypes.DenomUnit{{
						Denom:    res.GetNewTokenDenom(),
						Exponent: 0,
					}},
					Base:    res.GetNewTokenDenom(),
					Display: res.GetNewTokenDenom(),
					Name:    res.GetNewTokenDenom(),
					Symbol:  res.GetNewTokenDenom(),
				}, metadata)
			} else {
				s.Require().Error(err)
				// Ensure we don't charge if we expect an error
				s.Require().True(preCreateBalance.IsEqual(postCreateBalance))
			}
		})
	}
}

func (s *KeeperTestSuite) TestGasConsume() {
	// It's hard to estimate exactly how much gas will be consumed when creating a
	// denom, because besides consuming the gas specified by the params, the keeper
	// also does a bunch of other things that consume gas.
	//
	// Rather, we test whether the gas consumed is within a range. Specifically,
	// the range [gasConsume, gasConsume + offset]. If the actual gas consumption
	// falls within the range for all test cases, we consider the test passed.
	//
	// In experience, the total amount of gas consumed should consume be ~30k more
	// than the set amount.
	const offset = 50000

	addr0 := sample.AccAddress()

	for _, tc := range []struct {
		desc       string
		gasConsume uint64
	}{
		{
			desc:       "gas consume zero",
			gasConsume: 0,
		},
		{
			desc:       "gas consume 1,000,000",
			gasConsume: 1_000_000,
		},
		{
			desc:       "gas consume 10,000,000",
			gasConsume: 10_000_000,
		},
		{
			desc:       "gas consume 25,000,000",
			gasConsume: 25_000_000,
		},
		{
			desc:       "gas consume 50,000,000",
			gasConsume: 50_000_000,
		},
		{
			desc:       "gas consume 200,000,000",
			gasConsume: 200_000_000,
		},
	} {
		s.SetupTest()
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			// set params with the gas consume amount
			err := s.app.TokenfactoryKeeper.SetParams(s.ctx, types.NewParams(nil, tc.gasConsume))
			s.Require().NoError(err)

			// amount of gas consumed prior to the denom creation
			gasConsumedBefore := s.ctx.GasMeter().GasConsumed()

			// create a denom
			_, err = s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(addr0.String(), "larry"))
			s.Require().NoError(err)

			// amount of gas consumed after the denom creation
			gasConsumedAfter := s.ctx.GasMeter().GasConsumed()

			// the amount of gas consumed must be within the range
			gasConsumed := gasConsumedAfter - gasConsumedBefore
			s.Require().Greater(gasConsumed, tc.gasConsume)
			s.Require().Less(gasConsumed, tc.gasConsume+offset)
		})
	}
}
