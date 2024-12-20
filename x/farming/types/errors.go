package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/farming module sentinel errors
var (
	ErrInternal                = sdkerrors.Register(ModuleName, 1100, "internal error")
	ErrInvalidTime             = sdkerrors.Register(ModuleName, 1101, "invalid time")
	ErrInvalidToken            = sdkerrors.Register(ModuleName, 1102, "invalid token")
	ErrInvalidStakingAmount    = sdkerrors.Register(ModuleName, 1103, "invalid staking amount")
	ErrStakingBalanceNotEnough = sdkerrors.Register(ModuleName, 1104, "staking balance not enough")
	ErrNoReward                = sdkerrors.Register(ModuleName, 1105, "no reward")
)

func WrapError(err error, description string) error {
	return sdkerrors.Wrap(err, description)
}

func WrapErrorf(err error, format string, args ...interface{}) error {
	return sdkerrors.Wrapf(err, format, args...)
}

func WrapInternalError(err error) error {
	return sdkerrors.Wrapf(ErrInternal, ErrInternal.Error(), err)
}
