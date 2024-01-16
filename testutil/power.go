package testutil

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TokensToConsensusPower(tokens BigInt) int64 {
	return sdk.TokensToConsensusPower(sdkmath.NewIntFromBigInt(tokens.v), sdk.DefaultPowerReduction)
}

func TokensFromConsensusPower(power int64) BigInt {
	return MakeBigIntFromString(sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction).String())
}
