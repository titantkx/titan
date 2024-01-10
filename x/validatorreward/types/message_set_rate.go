package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetRate = "set_rate"

var _ sdk.Msg = &MsgSetRate{}

func NewMsgSetRate(authority string, rate sdk.Dec) *MsgSetRate {
	return &MsgSetRate{
		Authority: authority,
		Rate:      rate,
	}
}

func (msg *MsgSetRate) Route() string {
	return RouterKey
}

func (msg *MsgSetRate) Type() string {
	return TypeMsgSetRate
}

func (msg *MsgSetRate) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgSetRate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetRate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	// msg.Rate must provided
	if msg.Rate.IsNil() {
		return errorsmod.Wrapf(ErrInvalidRate, "rate must provided")
	}

	// msg.Rate must between 0 and 1
	if msg.Rate.LT(sdk.ZeroDec()) || msg.Rate.GT(sdk.OneDec()) {
		return errorsmod.Wrapf(ErrInvalidRate, "rate must between 0 and 1")
	}

	return nil
}
