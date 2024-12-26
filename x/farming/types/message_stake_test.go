package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
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
				Amount: utils.NewCoins("1000tkx"),
			},
		},
		{
			name: "invalid address",
			msg: MsgStake{
				Sender: "invalid_address",
				Amount: utils.NewCoins("1000tkx"),
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
