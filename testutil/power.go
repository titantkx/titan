package testutil

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TokensToConsensusPower(tokens Int) int64 {
	return sdk.TokensToConsensusPower(sdkmath.NewIntFromBigInt(tokens.v), sdk.DefaultPowerReduction)
}

func TokensFromConsensusPower(power int64) Int {
	return MakeIntFromString(sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction).String())
}
