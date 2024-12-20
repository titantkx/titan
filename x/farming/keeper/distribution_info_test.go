package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/farming/keeper"
	"github.com/titantkx/titan/x/farming/types"
)

func createTestDistributionInfo(keeper *keeper.Keeper, ctx sdk.Context) types.DistributionInfo {
	item := types.DistributionInfo{}
	keeper.SetDistributionInfo(ctx, item)
	return item
}

func TestDistributionInfoGet(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	item := createTestDistributionInfo(keeper, ctx)
	rst, found := keeper.GetDistributionInfo(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}
