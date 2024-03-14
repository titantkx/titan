package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
)

func TestMsgUpdateClass_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateClass
		err  error
	}{
		{
			name: "valid msg",
			msg: MsgUpdateClass{
				Creator: sample.AccAddress(),
				Id:      sample.ClassId(),
			},
		},
		{
			name: "invalid creator address",
			msg: MsgUpdateClass{
				Creator: "invalid_address",
				Id:      sample.ClassId(),
			},
			err: ErrInvalidAddress,
		},
		{
			name: "invalid class id",
			msg: MsgUpdateClass{
				Creator: sample.AccAddress(),
				Id:      "invalid_class_id",
			},
			err: ErrInvalidClassId,
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
