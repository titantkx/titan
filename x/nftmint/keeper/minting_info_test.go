package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/tokenize-titan/titan/testutil/keeper"
	"github.com/tokenize-titan/titan/testutil/nullify"
	"github.com/tokenize-titan/titan/x/nftmint/keeper"
	"github.com/tokenize-titan/titan/x/nftmint/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNMintingInfo(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.MintingInfo {
	items := make([]types.MintingInfo, n)
	for i := range items {
		items[i].ClassId = strconv.Itoa(i)

		keeper.SetMintingInfo(ctx, items[i])
	}
	return items
}

func TestMintingInfoGet(t *testing.T) {
	keeper, ctx := keepertest.NftmintKeeper(t)
	items := createNMintingInfo(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetMintingInfo(ctx,
			item.ClassId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestMintingInfoRemove(t *testing.T) {
	keeper, ctx := keepertest.NftmintKeeper(t)
	items := createNMintingInfo(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveMintingInfo(ctx,
			item.ClassId,
		)
		_, found := keeper.GetMintingInfo(ctx,
			item.ClassId,
		)
		require.False(t, found)
	}
}

func TestMintingInfoGetAll(t *testing.T) {
	keeper, ctx := keepertest.NftmintKeeper(t)
	items := createNMintingInfo(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllMintingInfo(ctx)),
	)
}
