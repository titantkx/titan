package keeper_test

import (
	"strconv"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/farming/keeper"
	"github.com/titantkx/titan/x/farming/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNStakingInfo(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.StakingInfo {
	items := make([]types.StakingInfo, n)
	for i := range items {
		items[i].Token = strconv.Itoa(i)
		items[i].Staker = strconv.Itoa(i)
		items[i].Amount = math.NewInt(int64(i))

		keeper.SetStakingInfo(ctx, items[i])
	}
	return items
}

func TestStakingInfoGet(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	items := createNStakingInfo(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetStakingInfo(ctx,
			item.Token,
			item.Staker,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestStakingInfoRemove(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	items := createNStakingInfo(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveStakingInfo(ctx,
			item.Token,
			item.Staker,
		)
		_, found := keeper.GetStakingInfo(ctx,
			item.Token,
			item.Staker,
		)
		require.False(t, found)
	}
}

func TestStakingInfoGetAll(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	items := createNStakingInfo(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllStakingInfo(ctx, "")),
	)
}
