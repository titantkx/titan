package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateDenom = "create_denom"

var _ sdk.Msg = &MsgCreateDenom{}

// NewMsgCreateDenom creates a msg to create a new denom
func NewMsgCreateDenom(sender, subdenom string) *MsgCreateDenom {
	return &MsgCreateDenom{
		Sender:   sender,
		Subdenom: subdenom,
	}
}

func (m MsgCreateDenom) Route() string { return RouterKey }

func (m MsgCreateDenom) Type() string { return TypeMsgCreateDenom }

func (m MsgCreateDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	_, err = GetTokenDenom(m.Sender, m.Subdenom)
	if err != nil {
		return errorsmod.Wrap(ErrInvalidDenom, err.Error())
	}

	return nil
}

func (m MsgCreateDenom) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}

func (m MsgCreateDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}
