package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/titantkx/titan/testutil/sample"
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
				Authority: "invalid_address",
				Rate:      sdk.NewDec(0),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSetRate{
				Authority: sample.AccAddress().String(),
				Rate:      sdk.NewDec(0),
			},
		}, {
			name: "rate not provided",
			msg: MsgSetRate{
				Authority: sample.AccAddress().String(),
			},
			err: ErrInvalidRate,
		}, {
			name: "rate too low",
			msg: MsgSetRate{
				Authority: sample.AccAddress().String(),
				Rate:      sdk.NewDec(-1),
			},
			err: ErrInvalidRate,
		}, {
			name: "rate too high",
			msg: MsgSetRate{
				Authority: sample.AccAddress().String(),
				Rate:      sdk.NewDec(1000000000000000000),
			},
			err: ErrInvalidRate,
		}, {
			name: "valid rate",
			msg: MsgSetRate{
				Authority: sample.AccAddress().String(),
				Rate:      sdk.NewDecWithPrec(1, 1),
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
