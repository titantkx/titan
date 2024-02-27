package cmd_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tokenize-titan/titan/testutil/cmd/bank"
	"github.com/tokenize-titan/titan/testutil/cmd/distribution"
	"github.com/tokenize-titan/titan/testutil/cmd/slashing"
	"github.com/tokenize-titan/titan/testutil/cmd/staking"
	"github.com/tokenize-titan/titan/testutil/cmd/status"
	"github.com/tokenize-titan/titan/utils"
)

func TestValidatorInactive(t *testing.T) {
	params := slashing.MustGetParams(t)

	// Create validator
	_, valBefore := MustCreateValidator(t, "")

	// Wait for one block to collect previous transaction fee's tax
	status.MustWait(t, status.MustGetLatestBlockHeight(t)+1)

	totalBalBefore := bank.MustGetTotalBalance(t, utils.BaseDenom, 0)
	distPoolBefore := distribution.MustGetCommunityPool(t).GetBaseDenomAmount()

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

	totalBalAfter := bank.MustGetTotalBalance(t, utils.BaseDenom, 0)
	distPoolAfter := distribution.MustGetCommunityPool(t).GetBaseDenomAmount()

	require.Equal(t, totalBalBefore, totalBalAfter)
	require.Equal(t, distPoolBefore.Add(slashedAmount), distPoolAfter)
}
