package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil/sample"
)

func TestMsgTransferClass_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgTransferClass
		err  error
	}{
		{
			name: "valid msg",
			msg: MsgTransferClass{
				Creator:  sample.AccAddress(),
				ClassId:  sample.ClassId(),
				Receiver: sample.AccAddress(),
			},
		},
		{
			name: "invalid creator address",
			msg: MsgTransferClass{
				Creator:  "invalid_address",
				ClassId:  sample.ClassId(),
				Receiver: sample.AccAddress(),
			},
			err: ErrInvalidAddress,
		},
		{
			name: "invalid receiver address",
			msg: MsgTransferClass{
				Creator:  sample.AccAddress(),
				ClassId:  sample.ClassId(),
				Receiver: "invalid_address",
			},
			err: ErrInvalidAddress,
		},
		{
			name: "invalid class id",
			msg: MsgTransferClass{
				Creator:  sample.AccAddress(),
				ClassId:  "invalid_class_id",
				Receiver: sample.AccAddress(),
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
