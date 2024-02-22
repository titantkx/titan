package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgTransferClass = "transfer_class"

var _ sdk.Msg = &MsgTransferClass{}

func NewMsgTransferClass(creator string, classId string, receiver string) *MsgTransferClass {
	return &MsgTransferClass{
		Creator:  creator,
		ClassId:  classId,
		Receiver: receiver,
	}
}

func (msg *MsgTransferClass) Route() string {
	return RouterKey
}

func (msg *MsgTransferClass) Type() string {
	return TypeMsgTransferClass
}

func (msg *MsgTransferClass) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTransferClass) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferClass) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return WrapErrorf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = strconv.ParseUint(msg.ClassId, 10, 64)
	if err != nil {
		return WrapErrorf(ErrInvalidClassId, "invalid class id (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return WrapErrorf(ErrInvalidAddress, "invalid receiver address (%s)", err)
	}

	return nil
}
