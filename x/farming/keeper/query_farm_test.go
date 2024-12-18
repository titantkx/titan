//nolint:dupl
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
	"github.com/titantkx/titan/x/farming/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestFarmQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNFarm(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryFarmRequest
		response *types.QueryFarmResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryFarmRequest{
				Token: msgs[0].Token,
			},
			response: &types.QueryFarmResponse{Farm: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryFarmRequest{
				Token: msgs[1].Token,
			},
			response: &types.QueryFarmResponse{Farm: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryFarmRequest{
				Token: strconv.Itoa(100000),
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
			response, err := keeper.Farm(wctx, tc.request)
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

func TestFarmQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNFarm(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryFarmAllRequest {
		return &types.QueryFarmAllRequest{
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
			//nolint:gosec // G115
			resp, err := keeper.FarmAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Farm), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Farm),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			//nolint:gosec // G115
			resp, err := keeper.FarmAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Farm), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Farm),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.FarmAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		//nolint:gosec // G115
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Farm),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.FarmAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
