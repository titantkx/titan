package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgHarvest = "harvest"

var _ sdk.Msg = &MsgHarvest{}

func NewMsgHarvest(sender string) *MsgHarvest {
	return &MsgHarvest{
		Sender: sender,
	}
}

func (msg *MsgHarvest) Route() string {
	return RouterKey
}

func (msg *MsgHarvest) Type() string {
	return TypeMsgHarvest
}

func (msg *MsgHarvest) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgHarvest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgHarvest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return WrapErrorf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}
