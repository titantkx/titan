//nolint:dupl
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUnstake = "unstake"

var _ sdk.Msg = &MsgUnstake{}

func NewMsgUnstake(sender string, amount sdk.Coins) *MsgUnstake {
	return &MsgUnstake{
		Sender: sender,
		Amount: amount,
	}
}

func (msg *MsgUnstake) Route() string {
	return RouterKey
}

func (msg *MsgUnstake) Type() string {
	return TypeMsgUnstake
}

func (msg *MsgUnstake) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgUnstake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnstake) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return WrapErrorf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.IsAllPositive() {
		return WrapError(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}
