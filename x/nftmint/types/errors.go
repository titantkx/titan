package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/nftmint module sentinel errors
var (
	ErrInternal       = sdkerrors.Register(ModuleName, 1100, "internal error: %s")
	ErrNotFound       = sdkerrors.Register(ModuleName, 1101, "not found")
	ErrUnauthorized   = sdkerrors.Register(ModuleName, 1102, "unauthorized")
	ErrInvalidAddress = sdkerrors.Register(ModuleName, 1103, "invalid address")
	ErrInvalidClassId = sdkerrors.Register(ModuleName, 1104, "invalid class id")
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
