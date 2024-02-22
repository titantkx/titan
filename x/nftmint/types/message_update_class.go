package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateClass = "update_class"

var _ sdk.Msg = &MsgUpdateClass{}

func NewMsgUpdateClass(creator string, id string, name string, symbol string, description string, uri string, uriHash string, data string) *MsgUpdateClass {
	return &MsgUpdateClass{
		Creator:     creator,
		Id:          id,
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Data:        data,
	}
}

func (msg *MsgUpdateClass) Route() string {
	return RouterKey
}

func (msg *MsgUpdateClass) Type() string {
	return TypeMsgUpdateClass
}

func (msg *MsgUpdateClass) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateClass) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateClass) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return WrapErrorf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = strconv.ParseUint(msg.Id, 10, 64)
	if err != nil {
		return WrapErrorf(ErrInvalidClassId, "invalid class id (%s)", err)
	}

	return nil
}
