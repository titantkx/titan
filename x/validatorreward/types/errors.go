package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/validatorreward module sentinel errors
var (
	ErrSample = errorsmod.Register(ModuleName, 1100, "sample error")

	ErrInvalidRate = errorsmod.Register(ModuleName, 2, "invalid rate")
	ErrForbidden   = errorsmod.Register(ModuleName, 3, "forbidden")
)
