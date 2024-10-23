package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

// TestMsgUpdateParams tests if valid/invalid update params messages are properly validated/invalidated
func TestMsgUpdateParams(t *testing.T) {
	coin := sdk.NewCoin("denom", math.NewInt(1))

	tests := []struct {
		name string
		msg  types.MsgUpdateParams
		err  error
	}{
		{
			name: "valid msg",
			msg: types.MsgUpdateParams{
				Authority: sample.AccAddress().String(),
				Params: types.Params{
					DenomCreationFee:        sdk.NewCoins(coin),
					DenomCreationGasConsume: 1000,
				},
			},
			err: nil,
		},
		{
			name: "invalid authority address",
			msg: types.MsgUpdateParams{
				Authority: "invalid address",
				Params: types.Params{
					DenomCreationFee:        sdk.NewCoins(coin),
					DenomCreationGasConsume: 1000,
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
