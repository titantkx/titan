package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddReward = "add_reward"

var _ sdk.Msg = &MsgAddReward{}

func NewMsgAddReward(sender string, token string, amount sdk.Coins, endTime time.Time, startTime time.Time) *MsgAddReward {
	return &MsgAddReward{
		Sender:    sender,
		Token:     token,
		Amount:    amount,
		EndTime:   endTime,
		StartTime: startTime,
	}
}

func (msg *MsgAddReward) Route() string {
	return RouterKey
}

func (msg *MsgAddReward) Type() string {
	return TypeMsgAddReward
}

func (msg *MsgAddReward) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgAddReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return WrapErrorf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if sdk.ValidateDenom(msg.Token) != nil {
		return WrapErrorf(ErrInvalidToken, "invalid token: %s", msg.Token)
	}

	if !msg.Amount.IsValid() || !msg.Amount.IsAllPositive() {
		return WrapError(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if msg.EndTime.IsZero() {
		return WrapError(ErrInvalidTime, "end time cannot be zero")
	}

	if !msg.StartTime.IsZero() && !msg.StartTime.Before(msg.EndTime) {
		return WrapError(ErrInvalidTime, "start time must be smaller than end time")
	}

	return nil
}
