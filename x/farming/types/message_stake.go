//nolint:dupl
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgStake = "stake"

var _ sdk.Msg = &MsgStake{}

func NewMsgStake(sender string, amount sdk.Coins) *MsgStake {
	return &MsgStake{
		Sender: sender,
		Amount: amount,
	}
}

func (msg *MsgStake) Route() string {
	return RouterKey
}

func (msg *MsgStake) Type() string {
	return TypeMsgStake
}

func (msg *MsgStake) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgStake) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return WrapErrorf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.IsAllPositive() {
		return WrapError(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}
