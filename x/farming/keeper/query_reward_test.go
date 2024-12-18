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

func TestRewardQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNReward(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryRewardRequest
		response *types.QueryRewardResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryRewardRequest{
				Farmer: msgs[0].Farmer,
			},
			response: &types.QueryRewardResponse{Reward: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryRewardRequest{
				Farmer: msgs[1].Farmer,
			},
			response: &types.QueryRewardResponse{Reward: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryRewardRequest{
				Farmer: strconv.Itoa(100000),
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
			response, err := keeper.Reward(wctx, tc.request)
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

func TestRewardQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNReward(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryRewardAllRequest {
		return &types.QueryRewardAllRequest{
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
			resp, err := keeper.RewardAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Reward), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Reward),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			//nolint:gosec // G115
			resp, err := keeper.RewardAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Reward), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Reward),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.RewardAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		//nolint:gosec // G115
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Reward),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.RewardAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
