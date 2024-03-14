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
				NewAuthority: sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSetAuthority{
				Authority:    sample.AccAddress(),
				NewAuthority: sample.AccAddress(),
			},
		}, {
			name: "invalid new address",
			msg: MsgSetAuthority{
				Authority:    sample.AccAddress(),
				NewAuthority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid new address",
			msg: MsgSetAuthority{
				Authority:    sample.AccAddress(),
				NewAuthority: sample.AccAddress(),
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
