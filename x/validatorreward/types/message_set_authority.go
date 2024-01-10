package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetAuthority = "set_authority"

var _ sdk.Msg = &MsgSetAuthority{}

func NewMsgSetAuthority(authority string, newAuthority string) *MsgSetAuthority {
	return &MsgSetAuthority{
		Authority:    authority,
		NewAuthority: newAuthority,
	}
}

func (msg *MsgSetAuthority) Route() string {
	return RouterKey
}

func (msg *MsgSetAuthority) Type() string {
	return TypeMsgSetAuthority
}

func (msg *MsgSetAuthority) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgSetAuthority) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetAuthority) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.NewAuthority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new authority address (%s)", err)
	}
	return nil
}
