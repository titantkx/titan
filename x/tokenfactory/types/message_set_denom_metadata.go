package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const TypeMsgSetDenomMetadata = "set_denom_metadata"

var _ sdk.Msg = &MsgSetDenomMetadata{}

// NewMsgChangeAdmin creates a message to burn tokens
func NewMsgSetDenomMetadata(sender string, metadata banktypes.Metadata) *MsgSetDenomMetadata {
	return &MsgSetDenomMetadata{
		Sender:   sender,
		Metadata: metadata,
	}
}

func (m MsgSetDenomMetadata) Route() string { return RouterKey }

func (m MsgSetDenomMetadata) Type() string { return TypeMsgSetDenomMetadata }

func (m MsgSetDenomMetadata) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	err = m.Metadata.Validate()
	if err != nil {
		return err
	}

	_, _, err = DeconstructDenom(m.Metadata.Base)
	if err != nil {
		return err
	}

	return nil
}

func (m MsgSetDenomMetadata) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}
