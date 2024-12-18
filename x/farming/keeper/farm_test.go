//nolint:dupl
package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/farming/keeper"
	"github.com/titantkx/titan/x/farming/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNFarm(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Farm {
	items := make([]types.Farm, n)
	for i := range items {
		items[i].Token = strconv.Itoa(i)

		keeper.SetFarm(ctx, items[i])
	}
	return items
}

func TestFarmGet(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	items := createNFarm(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetFarm(ctx,
			item.Token,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestFarmRemove(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	items := createNFarm(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveFarm(ctx,
			item.Token,
		)
		_, found := keeper.GetFarm(ctx,
			item.Token,
		)
		require.False(t, found)
	}
}

func TestFarmGetAll(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	items := createNFarm(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllFarm(ctx)),
	)
}
