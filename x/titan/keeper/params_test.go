package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/tokenize-titan/titan/testutil/keeper"
	"github.com/tokenize-titan/titan/x/titan/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.TitanKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
