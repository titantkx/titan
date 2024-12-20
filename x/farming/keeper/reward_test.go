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

func createNReward(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Reward {
	items := make([]types.Reward, n)
	for i := range items {
		items[i].Farmer = strconv.Itoa(i)

		keeper.SetReward(ctx, items[i])
	}
	return items
}

func TestRewardGet(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	items := createNReward(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetReward(ctx,
			item.Farmer,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestRewardRemove(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	items := createNReward(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveReward(ctx,
			item.Farmer,
		)
		_, found := keeper.GetReward(ctx,
			item.Farmer,
		)
		require.False(t, found)
	}
}

func TestRewardGetAll(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	items := createNReward(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllReward(ctx)),
	)
}
