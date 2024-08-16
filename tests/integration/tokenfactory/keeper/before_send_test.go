package keeper_test

import (
	"fmt"
	"os"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

type SendMsgTestCase struct {
	desc       string
	msg        func(denom string) *banktypes.MsgSend
	expectPass bool
}

func (s *KeeperTestSuite) TestBeforeSendHook() {
	s.SkipIfWSL()

	addr0 := sample.AccAddress()
	addr1 := sample.AccAddress()

	for _, tc := range []struct {
		desc     string
		wasmFile string
		sendMsgs []SendMsgTestCase
	}{
		{
			desc:     "should not allow sending 100 amount of *any* denom",
			wasmFile: "./testdata/no100.wasm",
			sendMsgs: []SendMsgTestCase{
				{
					desc: "sending 1 of factory denom should not error",
					msg: func(factorydenom string) *banktypes.MsgSend {
						return banktypes.NewMsgSend(
							addr0,
							addr1,
							sdk.NewCoins(sdk.NewInt64Coin(factorydenom, 1)),
						)
					},
					expectPass: true,
				},
				{
					desc: "sending 1 of non-factory denom should not error",
					msg: func(_ string) *banktypes.MsgSend {
						return banktypes.NewMsgSend(
							addr0,
							addr1,
							sdk.NewCoins(sdk.NewInt64Coin(utils.BaseDenom, 1)),
						)
					},
					expectPass: true,
				},
				{
					desc: "sending 100 of factory denom should error",
					msg: func(factorydenom string) *banktypes.MsgSend {
						return banktypes.NewMsgSend(
							addr0,
							addr1,
							sdk.NewCoins(sdk.NewInt64Coin(factorydenom, 100)),
						)
					},
					expectPass: false,
				},
				{
					desc: "sending 100 of non-factory denom should work",
					msg: func(_ string) *banktypes.MsgSend {
						return banktypes.NewMsgSend(
							addr0,
							addr1,
							sdk.NewCoins(sdk.NewInt64Coin(utils.BaseDenom, 100)),
						)
					},
					expectPass: true,
				},
				{
					desc: "having 100 coin within coins should not work",
					msg: func(factorydenom string) *banktypes.MsgSend {
						return banktypes.NewMsgSend(
							addr0,
							addr1,
							sdk.NewCoins(sdk.NewInt64Coin(factorydenom, 100), sdk.NewInt64Coin(utils.BaseDenom, 1)),
						)
					},
					expectPass: false,
				},
			},
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			// setup test
			s.SetupTest()

			// upload and instantiate wasm code
			wasmCode, err := os.ReadFile(tc.wasmFile)
			s.Require().NoError(err, "test: %v", tc.desc)
			codeID, _, err := s.contractKeeper.Create(s.ctx, addr0, wasmCode, nil)
			s.Require().NoError(err, "test: %v", tc.desc)
			cosmwasmAddress, _, err := s.contractKeeper.Instantiate(s.ctx, codeID, addr0, addr0, []byte("{}"), "", sdk.NewCoins())
			s.Require().NoError(err, "test: %v", tc.desc)

			// create new denom
			res, err := s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(addr0.String(), "bitcoin"))
			s.Require().NoError(err, "test: %v", tc.desc)
			denom := res.GetNewTokenDenom()

			// mint enough coins to the creator
			_, err = s.msgServer.Mint(s.ctx, types.NewMsgMint(addr0.String(), sdk.NewInt64Coin(denom, 1000000000)))
			s.Require().NoError(err)
			// mint some non token factory denom coins for testing
			s.FundAcc(addr0, sdk.Coins{sdk.NewInt64Coin(utils.BaseDenom, 100000000000)})

			// set beforesend hook to the new denom
			_, err = s.msgServer.SetBeforeSendHook(s.ctx, types.NewMsgSetBeforeSendHook(addr0.String(), denom, cosmwasmAddress.String()))
			s.Require().NoError(err, "test: %v", tc.desc)

			denoms, beforeSendHooks := s.app.TokenfactoryKeeper.GetAllBeforeSendHooks(s.ctx)
			s.Require().Equal(beforeSendHooks, []string{cosmwasmAddress.String()})
			s.Require().Equal(denoms, []string{denom})

			for _, sendTc := range tc.sendMsgs {
				_, err := s.bankMsgServer.Send(s.ctx, sendTc.msg(denom))
				if sendTc.expectPass {
					s.Require().NoError(err, "test: %v", sendTc.desc)
				} else {
					s.Require().Error(err, "test: %v", sendTc.desc)
				}

				// this is a check to ensure bank keeper wired in token factory keeper has hooks properly set
				// to check this, we try triggering bank hooks via token factory keeper
				for _, coin := range sendTc.msg(denom).Amount {
					_, err = s.msgServer.Mint(s.ctx, types.NewMsgMint(addr0.String(), sdk.NewInt64Coin(coin.Denom, coin.Amount.Int64())))
					if coin.Denom == denom && coin.Amount.Equal(math.NewInt(100)) {
						s.Require().Error(err, "test: %v", sendTc.desc)
					}
				}

			}
		})
	}
}

// func (s *KeeperTestSuite) TestInvalidSetBeforeSendHook() {
// 	s.SkipIfWSL()
// 	for _, tc := range []struct {
// 		desc     string
// 		wasmFile string
// 	}{
// 		{
// 			desc:     "should not allow sending 100 amount of *any* denom",
// 			wasmFile: "./testdata/no100.wasm",
// 		},
// 	} {
// 		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
// 			// setup test
// 			s.SetupTest()

// 			// upload and instantiate wasm code
// 			wasmCode, err := os.ReadFile(tc.wasmFile)
// 			s.Require().NoError(err, "test: %v", tc.desc)
// 			codeID, _, err := s.contractKeeper.Create(s.ctx, addr0, wasmCode, nil)
// 			s.Require().NoError(err, "test: %v", tc.desc)
// 			_, _, err = s.contractKeeper.Instantiate(s.ctx, codeID, addr0, addr0, []byte("{}"), "", sdk.NewCoins())
// 			s.Require().NoError(err, "test: %v", tc.desc)

// 			// create new denom
// 			res, err := s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(addr0.String(), "bitcoin"))
// 			s.Require().NoError(err, "test: %v", tc.desc)
// 			denom := res.GetNewTokenDenom()

// 			// mint enough coins to the creator
// 			_, err = s.msgServer.Mint(s.ctx, types.NewMsgMint(addr0.String(), sdk.NewInt64Coin(denom, 1000000000)))
// 			s.Require().NoError(err)
// 			// mint some non token factory denom coins for testing
// 			s.FundAcc(sdk.MustAccAddressFromBech32(addr0.String()), sdk.Coins{sdk.NewInt64Coin("foo", 100000000000)})

// 			// set an invalid beforesend hook to the new denom
// 			invalidAddress := sdk.AccAddress([]byte("addr1---------------"))
// 			_, err = s.msgServer.SetBeforeSendHook(s.ctx, types.NewMsgSetBeforeSendHook(addr0.String(), denom, invalidAddress.String()))
// 			s.Require().Error(err)

// 			denoms, beforeSendHooks := s.app.TokenFactoryKeeper.GetAllBeforeSendHooks(s.ctx)
// 			s.Require().Equal(beforeSendHooks, []string{})
// 			s.Require().Equal(denoms, []string{})
// 		})
// 	}
// }

// // TestInfiniteTrackBeforeSend tests gas metering with infinite loop contract
// // to properly test if we are gas metering trackBeforeSend properly.
// func (s *KeeperTestSuite) TestInfiniteTrackBeforeSend() {
// 	s.SkipIfWSL()

// 	for _, tc := range []struct {
// 		name            string
// 		wasmFile        string
// 		tokenToSend     sdk.Coins
// 		useFactoryDenom bool
// 		blockBeforeSend bool
// 		expectedError   bool
// 	}{
// 		{
// 			name:            "sending tokenfactory denom from module to module with infinite contract should panic",
// 			wasmFile:        "./testdata/infinite_track_beforesend.wasm",
// 			useFactoryDenom: true,
// 		},
// 		{
// 			name:            "sending tokenfactory denom from module to module with infinite contract should panic",
// 			wasmFile:        "./testdata/infinite_track_beforesend.wasm",
// 			useFactoryDenom: true,
// 			blockBeforeSend: true,
// 		},
// 		{
// 			name:            "sending non-tokenfactory denom from module to module with infinite contract should not panic",
// 			wasmFile:        "./testdata/infinite_track_beforesend.wasm",
// 			tokenToSend:     sdk.NewCoins(sdk.NewInt64Coin("foo", 1000000)),
// 			useFactoryDenom: false,
// 		},
// 		{
// 			name:            "Try using no 100 ",
// 			wasmFile:        "./testdata/no100.wasm",
// 			useFactoryDenom: true,
// 		},
// 	} {
// 		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
// 			// setup test
// 			s.SetupTest()

// 			// load wasm file
// 			wasmCode, err := os.ReadFile(tc.wasmFile)
// 			s.Require().NoError(err)

// 			// instantiate wasm code
// 			codeID, _, err := s.contractKeeper.Create(s.ctx, addr0, wasmCode, nil)
// 			s.Require().NoError(err, "test: %v", tc.name)
// 			cosmwasmAddress, _, err := s.contractKeeper.Instantiate(s.ctx, codeID, addr0, addr0, []byte("{}"), "", sdk.NewCoins())
// 			s.Require().NoError(err, "test: %v", tc.name)

// 			// create new denom
// 			res, err := s.msgServer.CreateDenom(s.ctx, types.NewMsgCreateDenom(addr0.String(), "bitcoin"))
// 			s.Require().NoError(err, "test: %v", tc.name)
// 			factoryDenom := res.GetNewTokenDenom()

// 			var tokenToSend sdk.Coins
// 			if tc.useFactoryDenom {
// 				tokenToSend = sdk.NewCoins(sdk.NewInt64Coin(factoryDenom, 100))
// 			} else {
// 				tokenToSend = tc.tokenToSend
// 			}

// 			// send the mint module tokenToSend
// 			if tc.blockBeforeSend {
// 				s.FundAcc(addr0, tokenToSend)
// 			} else {
// 				s.FundModuleAcc("mint", tokenToSend)
// 			}

// 			// set beforesend hook to the new denom
// 			// we register infinite loop contract here to test if we are gas metering properly
// 			_, err = s.msgServer.SetBeforeSendHook(s.ctx, types.NewMsgSetBeforeSendHook(addr0.String(), factoryDenom, cosmwasmAddress.String()))
// 			s.Require().NoError(err, "test: %v", tc.name)

// 			if tc.blockBeforeSend {
// 				err = s.app.BankKeeper.SendCoins(s.ctx, addr0, addr1, tokenToSend)
// 				s.Require().Error(err)
// 			} else {
// 				// track before send suppresses in any case, thus we expect no error
// 				err = s.app.BankKeeper.SendCoinsFromModuleToModule(s.ctx, "mint", "distribution", tokenToSend)
// 				s.Require().NoError(err)

// 				// send should happen regardless of trackBeforeSend results
// 				distributionModuleAddress := s.app.AccountKeeper.GetModuleAddress("distribution")
// 				distributionModuleBalances := s.app.BankKeeper.GetAllBalances(s.ctx, distributionModuleAddress)
// 				s.Require().True(distributionModuleBalances.Equal(tokenToSend))
// 			}

// 		})
// 	}
// }
