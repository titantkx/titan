package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	etherminttests "github.com/titantkx/ethermint/tests"

	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/pointer/keeper"
	"github.com/titantkx/titan/x/pointer/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNErc20Native(t *testing.T, keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Erc20Native {
	items := make([]types.Erc20Native, n)
	for i := range items {
		items[i].TokenDenom = strconv.Itoa(i)
		items[i].Erc20Addr = etherminttests.GenerateAddress().String()

		err := keeper.SetErc20Native(ctx, items[i])
		require.NoError(t, err)
	}

	return items
}

func TestSetErc20Native(t *testing.T) {
	testCases := []struct {
		name   string
		data   types.Erc20Native
		expErr interface{} // Use interface{} to handle both bool and specific error types
	}{
		{
			name: "correct hex and checksum",
			data: types.Erc20Native{
				Erc20Addr:  "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7",
				TokenDenom: "token3",
			},
			expErr: false,
		},
		{
			name: "wrong hex",
			data: types.Erc20Native{
				Erc20Addr:  "abcd",
				TokenDenom: "token1",
			},
			expErr: true,
		},
		{
			name: "correct hex but wrong checksum",
			data: types.Erc20Native{
				Erc20Addr:  "0xaBcD1234567890abcdef1234567890abcdef1200",
				TokenDenom: "token2",
			},
			expErr: errortypes.ErrInvalidAddress,
		},
		{
			name: "Do not accept zero address",
			data: types.Erc20Native{
				Erc20Addr:  "0x0000000000000000000000000000000000000000",
				TokenDenom: "token4",
			},
			expErr: errortypes.ErrInvalidAddress,
		},
		{
			name: "Empty address",
			data: types.Erc20Native{
				Erc20Addr:  "",
				TokenDenom: "token5",
			},
			expErr: true,
		},
		{
			name: "Invalid hex prefix",
			data: types.Erc20Native{
				Erc20Addr:  "0Y775b87ef5D82ca211811C1a02CE0fE0CA3a455d7",
				TokenDenom: "token6",
			},
			expErr: true,
		},
		{
			name: "Too short address",
			data: types.Erc20Native{
				Erc20Addr:  "0x775b87ef",
				TokenDenom: "token7",
			},
			expErr: true,
		},
	}

	keeper, ctx := keepertest.PointerKeeper(t)
	for _, tc := range testCases {
		err := keeper.SetErc20Native(ctx, tc.data)

		if tc.expErr == false {
			require.NoError(t, err, "expected no error for: %s", tc.name)
		} else if specificErr, ok := tc.expErr.(error); ok {
			require.ErrorIs(t, err, specificErr, "expected specific error for: %s", tc.name)
		} else {
			require.Error(t, err, "expected error for: %s", tc.name)
		}
	}
}

func TestErc20NativeGet(t *testing.T) {
	keeper, ctx := keepertest.PointerKeeper(t)
	items := createNErc20Native(t, keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetErc20Native(ctx,
			item.TokenDenom,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestErc20NativeRemove(t *testing.T) {
	keeper, ctx := keepertest.PointerKeeper(t)
	items := createNErc20Native(t, keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveErc20Native(ctx,
			item.TokenDenom,
		)
		_, found := keeper.GetErc20Native(ctx,
			item.TokenDenom,
		)
		require.False(t, found)
	}
}

func TestErc20NativeGetAll(t *testing.T) {
	keeper, ctx := keepertest.PointerKeeper(t)
	items := createNErc20Native(t, keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllErc20Native(ctx)),
	)
}
