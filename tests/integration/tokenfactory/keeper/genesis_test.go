package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func (s *KeeperTestSuite) TestGenesis() {
	genesisState := types.GenesisState{
		FactoryDenoms: []types.GenesisDenom{
			{
				Denom: "factory/titan1tk9ahr5eann843v82z50z8mx6p7lw7dm22t5kl/bitcoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "titan1tk9ahr5eann843v82z50z8mx6p7lw7dm22t5kl",
				},
			},
			{
				Denom: "factory/titan1tk9ahr5eann843v82z50z8mx6p7lw7dm22t5kl/diff-admin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "titan10wgwwzf7eyn9a83tjd2v3dl48mmhfkfsq5w3v3",
				},
			},
			{
				Denom: "factory/titan1tk9ahr5eann843v82z50z8mx6p7lw7dm22t5kl/litecoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "titan1tk9ahr5eann843v82z50z8mx6p7lw7dm22t5kl",
				},
			},
		},
	}

	s.SetupTestForInitGenesis()
	app := s.app

	// Test both with bank denom metadata set, and not set.
	for i, denom := range genesisState.FactoryDenoms {
		// hacky, sets bank metadata to exist if i != 0, to cover both cases.
		if i != 0 {
			app.BankKeeper.SetDenomMetaData(s.ctx, banktypes.Metadata{
				DenomUnits: []*banktypes.DenomUnit{{
					Denom:    denom.GetDenom(),
					Exponent: 0,
				}},
				Base:    denom.GetDenom(),
				Display: denom.GetDenom(),
				Name:    denom.GetDenom(),
				Symbol:  denom.GetDenom(),
			})
		}
	}

	// check before initGenesis that the module account is nil
	tokenfactoryModuleAccount := app.AccountKeeper.GetAccount(s.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().Nil(tokenfactoryModuleAccount)

	err := app.TokenfactoryKeeper.SetParams(s.ctx, types.Params{DenomCreationFee: sdk.Coins{sdk.NewInt64Coin(utils.BaseDenom, 100)}})
	s.Require().NoError(err)
	app.TokenfactoryKeeper.InitGenesis(s.ctx, genesisState)

	// check that the module account is now initialized
	tokenfactoryModuleAccount = app.AccountKeeper.GetAccount(s.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().NotNil(tokenfactoryModuleAccount)

	exportedGenesis := app.TokenfactoryKeeper.ExportGenesis(s.ctx)
	s.Require().NotNil(exportedGenesis)
	s.Require().Equal(genesisState, *exportedGenesis)

	// verify that the exported bank genesis is valid
	err = app.BankKeeper.SetParams(s.ctx, banktypes.DefaultParams())
	s.Require().NoError(err)
	exportedBankGenesis := app.BankKeeper.ExportGenesis(s.ctx)
	s.Require().NoError(exportedBankGenesis.Validate())

	app.BankKeeper.InitGenesis(s.ctx, exportedBankGenesis)
	for i, denom := range genesisState.FactoryDenoms {
		s.Require().NotNil(app.BankKeeper)

		// hacky, check whether bank metadata is not replaced if i != 0, to cover both cases.
		if i != 0 {
			metadata, found := app.BankKeeper.GetDenomMetaData(s.ctx, denom.GetDenom())
			s.Require().True(found)
			s.Require().EqualValues(metadata, banktypes.Metadata{
				DenomUnits: []*banktypes.DenomUnit{{
					Denom:    denom.GetDenom(),
					Exponent: 0,
				}},
				Base:    denom.GetDenom(),
				Display: denom.GetDenom(),
				Name:    denom.GetDenom(),
				Symbol:  denom.GetDenom(),
			})
		}
	}
}
