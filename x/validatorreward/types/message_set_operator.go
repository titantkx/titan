package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetOperator = "set_operator"

var _ sdk.Msg = &MsgSetOperator{}

func NewMsgSetOperator(operator string, newOperator string) *MsgSetOperator {
	return &MsgSetOperator{
		Operator:    operator,
		NewOperator: newOperator,
	}
}

func (msg *MsgSetOperator) Route() string {
	return RouterKey
}

func (msg *MsgSetOperator) Type() string {
	return TypeMsgSetOperator
}

func (msg *MsgSetOperator) GetSigners() []sdk.AccAddress {
	operator, err := sdk.AccAddressFromBech32(msg.Operator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{operator}
}

func (msg *MsgSetOperator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetOperator) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Operator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid operator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.NewOperator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new operator address (%s)", err)
	}
	return nil
}
