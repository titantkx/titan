package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgMint = "mint"

var _ sdk.Msg = &MsgMint{}

func NewMsgMint(creator string, receiver string, classId string, uri string, uriHash string, data string) *MsgMint {
	return &MsgMint{
		Creator:  creator,
		Receiver: receiver,
		ClassId:  classId,
		Uri:      uri,
		UriHash:  uriHash,
		Data:     data,
	}
}

func (msg *MsgMint) Route() string {
	return RouterKey
}

func (msg *MsgMint) Type() string {
	return TypeMsgMint
}

func (msg *MsgMint) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return WrapErrorf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return WrapErrorf(ErrInvalidAddress, "invalid receiver address (%s)", err)
	}

	_, err = strconv.ParseUint(msg.ClassId, 10, 64)
	if err != nil {
		return WrapErrorf(ErrInvalidClassId, "invalid class id (%s)", err)
	}

	return nil
}
