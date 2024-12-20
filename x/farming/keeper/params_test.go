package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/x/farming/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.FarmingKeeper(t)
	params := types.DefaultParams()

	err := k.SetParams(ctx, params)

	require.NoError(t, err)
	require.EqualValues(t, params, k.GetParams(ctx))
}
