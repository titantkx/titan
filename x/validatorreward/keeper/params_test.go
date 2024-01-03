package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/tokenize-titan/titan/testutil/keeper"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.ValidatorrewardKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
	require.EqualValues(t, params.Rate, k.Rate(ctx))
	require.EqualValues(t, params.Operator, k.Operator(ctx))
}
