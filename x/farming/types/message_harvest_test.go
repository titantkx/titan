package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
)

func TestMsgHarvest_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgHarvest
		err  error
	}{
		{
			name: "valid address",
			msg: MsgHarvest{
				Sender: sample.AccAddress().String(),
			},
		},
		{
			name: "invalid address",
			msg: MsgHarvest{
				Sender: "invalid_address",
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
