package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgFundRewardPool = "fund_reward_pool"

var _ sdk.Msg = &MsgFundRewardPool{}

func NewMsgFundRewardPool(depositor sdk.AccAddress, amount sdk.Coins) *MsgFundRewardPool {
	return &MsgFundRewardPool{
		Depositor: depositor.String(),
		Amount:    amount,
	}
}

func (msg *MsgFundRewardPool) Route() string {
	return RouterKey
}

func (msg *MsgFundRewardPool) Type() string {
	return TypeMsgFundRewardPool
}

func (msg *MsgFundRewardPool) GetSigners() []sdk.AccAddress {
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{depositor}
}

func (msg *MsgFundRewardPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFundRewardPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid depositor address (%s)", err)
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid amount (%s)", msg.Amount.String())
	}
	return nil
}
