package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
)

func TestMsgSetAuthority_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSetAuthority
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSetAuthority{
				Authority:    "invalid_address",
				NewAuthority: sample.AccAddress().String(),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSetAuthority{
				Authority:    sample.AccAddress().String(),
				NewAuthority: sample.AccAddress().String(),
			},
		}, {
			name: "invalid new address",
			msg: MsgSetAuthority{
				Authority:    sample.AccAddress().String(),
				NewAuthority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid new address",
			msg: MsgSetAuthority{
				Authority:    sample.AccAddress().String(),
				NewAuthority: sample.AccAddress().String(),
			},
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
