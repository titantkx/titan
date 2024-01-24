package cmd_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tokenize-titan/titan/testutil/cmd/distribution"
	"github.com/tokenize-titan/titan/testutil/cmd/slashing"
	"github.com/tokenize-titan/titan/testutil/cmd/staking"
)

func TestValidatorInactive(t *testing.T) {
	params := slashing.MustGetParams(t)

	// Create validator
	_, valBefore := MustCreateValidator(t, "")

	totalBalBefore := MustGetTotalBalance(t, 0)
	distPoolBefore := distribution.MustGetCommunityPool(t)

	// Wait until validator is jailed for being inactive
	var valAfter staking.Validator
	for {
		time.Sleep(1 * time.Second)
		valAfter = staking.MustGetValidator(t, valBefore.OperatorAddress)
		if valAfter.Jailed {
			break
		}
	}

	slashedAmount := valBefore.Tokens.Float().Mul(params.SlashFractionDowntime).Int()

	require.Equal(t, valBefore.Tokens.Sub(slashedAmount), valAfter.Tokens)

	totalBalAfter := MustGetTotalBalance(t, 0)
	distPoolAfter := distribution.MustGetCommunityPool(t)

	require.Equal(t, totalBalBefore, totalBalAfter)
	require.Equal(t, distPoolBefore.Pool.GetBaseDenomAmount().Add(slashedAmount), distPoolAfter.Pool.GetBaseDenomAmount())
}
