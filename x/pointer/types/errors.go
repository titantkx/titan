package types

// DONTCOVER

import (
	// sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// x/pointer module sentinel errors
var (
	ErrInternal = sdkerrors.Register(ModuleName, 1100, "internal error")
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
