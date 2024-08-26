package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func (s *KeeperTestSuite) TestAdminMsgs() {
	addr0 := sample.AccAddress()
	addr1 := sample.AccAddress()

	addr0bal := int64(0)
	addr1bal := int64(0)

	bankKeeper := s.app.BankKeeper

	denom := s.CreateDenom(addr0.String(), "bitcoin")

	// Make sure that the admin is set correctly
	queryRes, err := s.queryClient.DenomAuthorityMetadata(s.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{Denom: denom})
	s.Require().NoError(err)
	s.Require().Equal(addr0.String(), queryRes.AuthorityMetadata.Admin)

	// Test minting to admins own account
	_, err = s.msgServer.Mint(s.ctx, types.NewMsgMint(addr0.String(), sdk.NewInt64Coin(denom, 10)))
	addr0bal += 10
	s.Require().NoError(err)
	s.Require().True(bankKeeper.GetBalance(s.ctx, addr0, denom).Amount.Int64() == addr0bal, bankKeeper.GetBalance(s.ctx, addr0, denom))

	// Test minting to a different account
	_, err = s.msgServer.Mint(s.ctx, types.NewMsgMintTo(addr0.String(), sdk.NewInt64Coin(denom, 10), addr1.String()))
	addr1bal += 10
	s.Require().NoError(err)
	s.Require().True(s.app.BankKeeper.GetBalance(s.ctx, addr1, denom).Amount.Int64() == addr1bal, s.app.BankKeeper.GetBalance(s.ctx, addr1, denom))

	// Test burning from own account
	_, err = s.msgServer.Burn(s.ctx, types.NewMsgBurn(addr0.String(), sdk.NewInt64Coin(denom, 5)))
	s.Require().NoError(err)
	s.Require().True(bankKeeper.GetBalance(s.ctx, addr1, denom).Amount.Int64() == addr1bal)

	// Test Change Admin
	_, err = s.msgServer.ChangeAdmin(s.ctx, types.NewMsgChangeAdmin(addr0.String(), denom, addr1.String()))
	s.Require().NoError(err)
	queryRes, err = s.queryClient.DenomAuthorityMetadata(s.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
		Denom: denom,
	})
	s.Require().NoError(err)
	s.Require().Equal(addr1.String(), queryRes.AuthorityMetadata.Admin)

	// Make sure old admin can no longer do actions
	_, err = s.msgServer.Burn(s.ctx, types.NewMsgBurn(addr0.String(), sdk.NewInt64Coin(denom, 5)))
	s.Require().Error(err)

	// Make sure the new admin works
	_, err = s.msgServer.Mint(s.ctx, types.NewMsgMint(addr1.String(), sdk.NewInt64Coin(denom, 5)))
	addr1bal += 5
	s.Require().NoError(err)
	s.Require().True(bankKeeper.GetBalance(s.ctx, addr1, denom).Amount.Int64() == addr1bal)

	// Try setting admin to empty
	_, err = s.msgServer.ChangeAdmin(s.ctx, types.NewMsgChangeAdmin(addr1.String(), denom, ""))
	s.Require().NoError(err)
	queryRes, err = s.queryClient.DenomAuthorityMetadata(s.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{Denom: denom})
	s.Require().NoError(err)
	s.Require().Equal("", queryRes.AuthorityMetadata.Admin)
}

// TestMintDenom ensures the following properties of the MintMessage:
// * No one can mint tokens for a denom that doesn't exist
// * Only the admin of a denom can mint tokens for it
// * The admin of a denom can mint tokens for it
func (s *KeeperTestSuite) TestMintDenom() {
	addr0 := sample.AccAddress()
	addr1 := sample.AccAddress()

	balances := map[string]int64{
		addr0.String(): 0,
		addr1.String(): 0,
	}

	denom := s.CreateDenom(addr0.String(), "bitcoin")

	for _, tc := range []struct {
		desc       string
		mintMsg    types.MsgMint
		expectPass bool
	}{
		{
			desc: "denom does not exist",
			mintMsg: *types.NewMsgMint(
				addr0.String(),
				sdk.NewInt64Coin("factory/titan1tk9ahr5eann843v82z50z8mx6p7lw7dm22t5kl/evmos", 10),
			),
			expectPass: false,
		},
		{
			desc: "mint is not by the admin",
			mintMsg: *types.NewMsgMintTo(
				addr1.String(),
				sdk.NewInt64Coin(denom, 10),
				addr0.String(),
			),
			expectPass: false,
		},
		{
			desc: "success case - mint to self",
			mintMsg: *types.NewMsgMint(
				addr0.String(),
				sdk.NewInt64Coin(denom, 10),
			),
			expectPass: true,
		},
		{
			desc: "success case - mint to another address",
			mintMsg: *types.NewMsgMintTo(
				addr0.String(),
				sdk.NewInt64Coin(denom, 10),
				addr1.String(),
			),
			expectPass: true,
		},
		{
			desc: "error: try minting non-tokenfactory denom",
			mintMsg: *types.NewMsgMintTo(
				addr0.String(),
				sdk.NewInt64Coin(utils.BaseDenom, 10),
				addr1.String(),
			),
			expectPass: false,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			_, err := s.msgServer.Mint(s.ctx, &tc.mintMsg)
			if tc.expectPass {
				s.Require().NoError(err)
				balances[tc.mintMsg.MintToAddress] += tc.mintMsg.Amount.Amount.Int64()
			} else {
				s.Require().Error(err)
			}

			mintToAddr, _ := sdk.AccAddressFromBech32(tc.mintMsg.MintToAddress)
			bal := s.app.BankKeeper.GetBalance(s.ctx, mintToAddr, denom).Amount
			s.Require().Equal(bal.Int64(), balances[tc.mintMsg.MintToAddress])
		})
	}
}

func (s *KeeperTestSuite) TestBurnDenom() {
	addr0 := sample.AccAddress()
	addr1 := sample.AccAddress()

	denom := s.CreateDenom(addr0.String(), "bitcoin")

	// mint 1000 tokens for all test accounts
	_, err := s.msgServer.Mint(s.ctx, types.NewMsgMintTo(addr0.String(), sdk.NewInt64Coin(denom, 1000), addr0.String()))
	s.Require().NoError(err)

	_, err = s.msgServer.Mint(s.ctx, types.NewMsgMintTo(addr0.String(), sdk.NewInt64Coin(denom, 1000), addr1.String()))
	s.Require().NoError(err)

	balances := map[string]int64{
		addr0.String(): 1000,
		addr1.String(): 1000,
	}

	// save sample module account address for testing
	moduleAdress := s.app.AccountKeeper.GetModuleAddress(types.ModuleName)

	for _, tc := range []struct {
		desc       string
		burnMsg    types.MsgBurn
		expectPass bool
	}{
		{
			desc: "denom does not exist",
			burnMsg: *types.NewMsgBurn(
				addr0.String(),
				sdk.NewInt64Coin("factory/titan1tk9ahr5eann843v82z50z8mx6p7lw7dm22t5kl/evmos", 10),
			),
			expectPass: false,
		},
		{
			desc: "burn is not by the admin",
			burnMsg: *types.NewMsgBurnFrom(
				addr1.String(),
				sdk.NewInt64Coin(denom, 10),
				addr0.String(),
			),
			expectPass: false,
		},
		{
			desc: "burn more than balance",
			burnMsg: *types.NewMsgBurn(
				addr0.String(),
				sdk.NewInt64Coin(denom, 10000),
			),
			expectPass: false,
		},
		{
			desc: "success case - burn from self",
			burnMsg: *types.NewMsgBurn(
				addr0.String(),
				sdk.NewInt64Coin(denom, 10),
			),
			expectPass: true,
		},
		{
			desc: "success case - burn from another address",
			burnMsg: *types.NewMsgBurnFrom(
				addr0.String(),
				sdk.NewInt64Coin(denom, 10),
				addr1.String(),
			),
			expectPass: true,
		},
		{
			desc: "fail case - burn from module account",
			burnMsg: *types.NewMsgBurnFrom(
				addr0.String(),
				sdk.NewInt64Coin(denom, 10),
				moduleAdress.String(),
			),
			expectPass: false,
		},
		{
			desc: "fail case - burn non-tokenfactory denom",
			burnMsg: *types.NewMsgBurnFrom(
				addr0.String(),
				sdk.NewInt64Coin(utils.BaseDenom, 10),
				moduleAdress.String(),
			),
			expectPass: false,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			_, err := s.msgServer.Burn(s.ctx, &tc.burnMsg)
			if tc.expectPass {
				s.Require().NoError(err)
				balances[tc.burnMsg.BurnFromAddress] -= tc.burnMsg.Amount.Amount.Int64()
			} else {
				s.Require().Error(err)
			}

			burnFromAddr, _ := sdk.AccAddressFromBech32(tc.burnMsg.BurnFromAddress)
			bal := s.app.BankKeeper.GetBalance(s.ctx, burnFromAddr, denom).Amount
			s.Require().Equal(bal.Int64(), balances[tc.burnMsg.BurnFromAddress])
		})
	}
}

func (s *KeeperTestSuite) TestChangeAdminDenom() {
	addr0 := sample.AccAddress()
	addr1 := sample.AccAddress()
	addr2 := sample.AccAddress()

	for _, tc := range []struct {
		desc                    string
		msgChangeAdmin          func(denom string) *types.MsgChangeAdmin
		expectedChangeAdminPass bool
		expectedAdmin           string
		msgMint                 func(denom string) *types.MsgMint
		expectedMintPass        bool
	}{
		{
			desc: "creator admin can't mint after setting to '' ",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(addr0.String(), denom, "")
			},
			expectedChangeAdminPass: true,
			expectedAdmin:           "",
			msgMint: func(denom string) *types.MsgMint {
				return types.NewMsgMint(addr0.String(), sdk.NewInt64Coin(denom, 5))
			},
			expectedMintPass: false,
		},
		{
			desc: "non-admins can't change the existing admin",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(addr1.String(), denom, addr2.String())
			},
			expectedChangeAdminPass: false,
			expectedAdmin:           addr0.String(),
		},
		{
			desc: "success change admin",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(addr0.String(), denom, addr1.String())
			},
			expectedAdmin:           addr1.String(),
			expectedChangeAdminPass: true,
			msgMint: func(denom string) *types.MsgMint {
				return types.NewMsgMint(addr1.String(), sdk.NewInt64Coin(denom, 5))
			},
			expectedMintPass: true,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			// setup test
			s.SetupTest()

			// Create a denom and mint
			res, err := s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(addr0.String(), "bitcoin"))
			s.Require().NoError(err)

			testDenom := res.GetNewTokenDenom()

			_, err = s.msgServer.Mint(s.ctx, types.NewMsgMint(addr0.String(), sdk.NewInt64Coin(testDenom, 10)))
			s.Require().NoError(err)

			_, err = s.msgServer.ChangeAdmin(s.ctx, tc.msgChangeAdmin(testDenom))
			if tc.expectedChangeAdminPass {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}

			queryRes, err := s.queryClient.DenomAuthorityMetadata(s.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
				Denom: testDenom,
			})
			s.Require().NoError(err)

			s.Require().Equal(tc.expectedAdmin, queryRes.AuthorityMetadata.Admin)

			// we test mint to test if admin authority is performed properly after admin change.
			if tc.msgMint != nil {
				_, err := s.msgServer.Mint(s.ctx, tc.msgMint(testDenom))
				if tc.expectedMintPass {
					s.Require().NoError(err)
				} else {
					s.Require().Error(err)
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestSetDenomMetaData() {
	addr0 := sample.AccAddress()
	addr1 := sample.AccAddress()

	denom := s.CreateDenom(addr0.String(), "bitcoin")

	for _, tc := range []struct {
		desc                string
		msgSetDenomMetadata types.MsgSetDenomMetadata
		expectedPass        bool
	}{
		{
			desc: "successful set denom metadata",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(addr0.String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    denom,
						Exponent: 0,
					},
					{
						Denom:    utils.BaseDenom,
						Exponent: 6,
					},
				},
				Base:    denom,
				Display: utils.BaseDenom,
				Name:    "TKX",
				Symbol:  "TKX",
			}),
			expectedPass: true,
		},
		{
			desc: "non existent factory denom name",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(addr0.String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    fmt.Sprintf("factory/%s/litecoin", addr0.String()),
						Exponent: 0,
					},
					{
						Denom:    utils.BaseDenom,
						Exponent: 6,
					},
				},
				Base:    fmt.Sprintf("factory/%s/litecoin", addr0.String()),
				Display: utils.BaseDenom,
				Name:    "TKX",
				Symbol:  "TKX",
			}),
			expectedPass: false,
		},
		{
			desc: "non-factory denom",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(addr0.String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    utils.BaseDenom,
						Exponent: 0,
					},
					{
						Denom:    "tkx",
						Exponent: 6,
					},
				},
				Base:    utils.BaseDenom,
				Display: "tkx",
				Name:    "TKX",
				Symbol:  "TKX",
			}),
			expectedPass: false,
		},
		{
			desc: "wrong admin",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(addr1.String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    denom,
						Exponent: 0,
					},
					{
						Denom:    utils.BaseDenom,
						Exponent: 6,
					},
				},
				Base:    denom,
				Display: utils.BaseDenom,
				Name:    "TKX",
				Symbol:  "TKX",
			}),
			expectedPass: false,
		},
		{
			desc: "invalid metadata (missing display denom unit)",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(addr0.String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    denom,
						Exponent: 0,
					},
				},
				Base:    denom,
				Display: utils.BaseDenom,
				Name:    "TKX",
				Symbol:  "TKX",
			}),
			expectedPass: false,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			bankKeeper := s.app.BankKeeper
			res, err := s.msgServer.SetDenomMetadata(s.ctx, &tc.msgSetDenomMetadata)
			if tc.expectedPass {
				s.Require().NoError(err)
				s.Require().NotNil(res)

				md, found := bankKeeper.GetDenomMetaData(s.ctx, denom)
				s.Require().True(found)
				s.Require().Equal(tc.msgSetDenomMetadata.Metadata.Name, md.Name)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
