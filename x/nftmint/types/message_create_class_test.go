package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
)

func TestMsgCreateClass_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateClass
		err  error
	}{
		{
			name: "valid msg",
			msg: MsgCreateClass{
				Creator: sample.AccAddress().String(),
			},
		},
		{
			name: "invalid creator address",
			msg: MsgCreateClass{
				Creator: "invalid_address",
			},
			err: ErrInvalidAddress,
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
