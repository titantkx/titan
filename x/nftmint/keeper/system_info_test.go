package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/nftmint/keeper"
	"github.com/titantkx/titan/x/nftmint/types"
)

func createTestSystemInfo(keeper *keeper.Keeper, ctx sdk.Context) types.SystemInfo {
	item := types.SystemInfo{}
	keeper.SetSystemInfo(ctx, item)
	return item
}

func TestSystemInfoGet(t *testing.T) {
	keeper, ctx := keepertest.NftmintKeeper(t)
	item := createTestSystemInfo(keeper, ctx)
	rst, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestSystemInfoRemove(t *testing.T) {
	keeper, ctx := keepertest.NftmintKeeper(t)
	createTestSystemInfo(keeper, ctx)
	keeper.RemoveSystemInfo(ctx)
	_, found := keeper.GetSystemInfo(ctx)
	require.False(t, found)
}
