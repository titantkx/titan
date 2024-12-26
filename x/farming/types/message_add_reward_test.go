package types

import (
	"testing"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
)

func TestMsgAddReward_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddReward
		err  error
	}{
		{
			name: "valid",
			msg: MsgAddReward{
				Sender:    sample.AccAddress().String(),
				Token:     "btc",
				Amount:    utils.NewCoins("1000tkx"),
				EndTime:   time.Now().Add(1 * time.Hour),
				StartTime: time.Now(),
			},
		},
		{
			name: "valid zero start time",
			msg: MsgAddReward{
				Sender:  sample.AccAddress().String(),
				Token:   "btc",
				Amount:  utils.NewCoins("1000tkx"),
				EndTime: time.Now().Add(1 * time.Hour),
			},
		},
		{
			name: "invalid token",
			msg: MsgAddReward{
				Sender:    sample.AccAddress().String(),
				Token:     "123",
				Amount:    utils.NewCoins("1000tkx"),
				EndTime:   time.Now().Add(1 * time.Hour),
				StartTime: time.Now(),
			},
			err: ErrInvalidToken,
		},
		{
			name: "invalid address",
			msg: MsgAddReward{
				Sender:    "invalid_address",
				Token:     "btc",
				Amount:    utils.NewCoins("1000tkx"),
				EndTime:   time.Now().Add(1 * time.Hour),
				StartTime: time.Now(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "zero reward",
			msg: MsgAddReward{
				Sender:    sample.AccAddress().String(),
				Token:     "btc",
				EndTime:   time.Now().Add(1 * time.Hour),
				StartTime: time.Now(),
			},
			err: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "zero end time",
			msg: MsgAddReward{
				Sender:    sample.AccAddress().String(),
				Token:     "btc",
				Amount:    utils.NewCoins("1000tkx"),
				StartTime: time.Now(),
			},
			err: ErrInvalidTime,
		},
		{
			name: "start after end",
			msg: MsgAddReward{
				Sender:    sample.AccAddress().String(),
				Token:     "btc",
				Amount:    utils.NewCoins("1000tkx"),
				EndTime:   time.Now().Add(1 * time.Hour),
				StartTime: time.Now().Add(2 * time.Hour),
			},
			err: ErrInvalidTime,
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
