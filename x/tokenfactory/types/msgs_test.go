package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztesting "github.com/titantkx/titan/testutil/authz"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/x/tokenfactory"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

// Test authz serialize and de-serializes for tokenfactory msg.
func TestAuthzMsg(t *testing.T) {
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address()).String()
	coin := sdk.NewCoin("denom", math.NewInt(1))

	testCases := []struct {
		name string
		msg  sdk.Msg
	}{
		{
			name: "MsgCreateDenom",
			msg: &types.MsgCreateDenom{
				Sender:   addr1,
				Subdenom: "valoper1xyz",
			},
		},
		{
			name: "MsgBurn",
			msg: &types.MsgBurn{
				Sender: addr1,
				Amount: coin,
			},
		},
		{
			name: "MsgMint",
			msg: &types.MsgMint{
				Sender: addr1,
				Amount: coin,
			},
		},
		{
			name: "MsgChangeAdmin",
			msg: &types.MsgChangeAdmin{
				Sender:   addr1,
				Denom:    "denom",
				NewAdmin: sample.AccAddress().String(),
			},
		},
		{
			name: "MsgUpdateParams",
			msg: &types.MsgUpdateParams{
				Authority: addr1,
				Params: types.Params{
					DenomCreationFee:        sdk.NewCoins(coin),
					DenomCreationGasConsume: 1000,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authztesting.TestMessageAuthzSerialization(t, tc.msg, tokenfactory.AppModuleBasic{})
		})
	}
}
