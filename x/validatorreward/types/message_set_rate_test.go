package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tokenize-titan/titan/testutil/sample"
)

func TestMsgSetRate_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSetRate
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSetRate{
				Operator: "invalid_address",
				Rate:     sdk.NewDec(0),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSetRate{
				Operator: sample.AccAddress(),
				Rate:     sdk.NewDec(0),
			},
		}, {
			name: "rate not provided",
			msg: MsgSetRate{
				Operator: sample.AccAddress(),
			},
			err: ErrInvalidRate,
		}, {
			name: "rate too low",
			msg: MsgSetRate{
				Operator: sample.AccAddress(),
				Rate:     sdk.NewDec(-1),
			},
			err: ErrInvalidRate,
		}, {
			name: "rate too high",
			msg: MsgSetRate{
				Operator: sample.AccAddress(),
				Rate:     sdk.NewDec(1000000000000000000),
			},
			err: ErrInvalidRate,
		}, {
			name: "valid rate",
			msg: MsgSetRate{
				Operator: sample.AccAddress(),
				Rate:     sdk.NewDec(0),
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
