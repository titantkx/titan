package types

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
)

func TestMsgMint_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgMint
		err  error
	}{
		{
			name: "valid msg",
			msg: MsgMint{
				Creator:  sample.AccAddress().String(),
				Receiver: sample.AccAddress().String(),
				ClassId:  sample.ClassId(),
			},
		},
		{
			name: "invalid creator address",
			msg: MsgMint{
				Creator:  "invalid_address",
				Receiver: sample.AccAddress().String(),
				ClassId:  strconv.FormatUint(rand.Uint64(), 10),
			},
			err: ErrInvalidAddress,
		},
		{
			name: "invalid receiver address",
			msg: MsgMint{
				Creator:  sample.AccAddress().String(),
				Receiver: "invalid_address",
				ClassId:  sample.ClassId(),
			},
			err: ErrInvalidAddress,
		},
		{
			name: "invalid class id",
			msg: MsgMint{
				Creator:  sample.AccAddress().String(),
				Receiver: sample.AccAddress().String(),
				ClassId:  "invalid_class_id",
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
