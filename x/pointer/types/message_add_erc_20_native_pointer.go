package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddERC20NativePointer = "add_erc_20_native_pointer"

var _ sdk.Msg = &MsgAddERC20NativePointer{}

func NewMsgAddERC20NativePointer(authority string, token string, name string, symbol string, decimals uint64) *MsgAddERC20NativePointer {
	return &MsgAddERC20NativePointer{
		Authority: authority,
		Token:     token,
		Name:      name,
		Symbol:    symbol,
		Decimals:  decimals,
	}
}

func (msg *MsgAddERC20NativePointer) Route() string {
	return RouterKey
}

func (msg *MsgAddERC20NativePointer) Type() string {
	return TypeMsgAddERC20NativePointer
}

func (msg *MsgAddERC20NativePointer) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgAddERC20NativePointer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddERC20NativePointer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return WrapErrorf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	// decimals must be between 0 and 18
	if msg.Decimals > 18 {
		return WrapError(sdkerrors.ErrInvalidRequest, "decimals must be between 0 and 18")
	}
	return nil
}
