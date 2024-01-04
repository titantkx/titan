package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetRate = "set_rate"

var _ sdk.Msg = &MsgSetRate{}

func NewMsgSetRate(operator string, rate sdk.Dec) *MsgSetRate {
	return &MsgSetRate{
		Operator: operator,
		Rate:     rate,
	}
}

func (msg *MsgSetRate) Route() string {
	return RouterKey
}

func (msg *MsgSetRate) Type() string {
	return TypeMsgSetRate
}

func (msg *MsgSetRate) GetSigners() []sdk.AccAddress {
	operator, err := sdk.AccAddressFromBech32(msg.Operator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{operator}
}

func (msg *MsgSetRate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetRate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Operator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid operator address (%s)", err)
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
