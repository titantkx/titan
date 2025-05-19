package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/pointer/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestErc20NativeQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.PointerKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNErc20Native(t, keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryErc20NativeRequest
		response *types.QueryErc20NativeResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryErc20NativeRequest{
				TokenDenom: msgs[0].TokenDenom,
			},
			response: &types.QueryErc20NativeResponse{Erc20Native: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryErc20NativeRequest{
				TokenDenom: msgs[1].TokenDenom,
			},
			response: &types.QueryErc20NativeResponse{Erc20Native: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryErc20NativeRequest{
				TokenDenom: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Erc20Native(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestErc20NativeQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.PointerKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNErc20Native(t, keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryErc20NativeAllRequest {
		return &types.QueryErc20NativeAllRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.Erc20NativeAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Erc20Native), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Erc20Native),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.Erc20NativeAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Erc20Native), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Erc20Native),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.Erc20NativeAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Erc20Native),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.Erc20NativeAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
