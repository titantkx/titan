package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil/sample"
)

func TestMsgSetOperator_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSetOperator
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSetOperator{
				Operator:    "invalid_address",
				NewOperator: sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSetOperator{
				Operator:    sample.AccAddress(),
				NewOperator: sample.AccAddress(),
			},
		}, {
			name: "invalid new address",
			msg: MsgSetOperator{
				Operator:    sample.AccAddress(),
				NewOperator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid new address",
			msg: MsgSetOperator{
				Operator:    sample.AccAddress(),
				NewOperator: sample.AccAddress(),
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
