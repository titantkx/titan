package types

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
)

func TestMsgStake_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgStake
		err  error
	}{
		{
			name: "valid",
			msg: MsgStake{
				Sender: sample.AccAddress().String(),
				Amount: sdk.NewCoins(sdk.NewCoin("tkx", math.NewInt(500000000))),
			},
		},
		{
			name: "invalid address",
			msg: MsgStake{
				Sender: "invalid_address",
				Amount: sdk.NewCoins(sdk.NewCoin("tkx", math.NewInt(500000000))),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "zero amount",
			msg: MsgStake{
				Sender: sample.AccAddress().String(),
			},
			err: sdkerrors.ErrInvalidCoins,
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
