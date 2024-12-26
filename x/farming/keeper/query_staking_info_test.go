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

func TestStakingInfoQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNStakingInfo(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryStakingInfoRequest
		response *types.QueryStakingInfoResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryStakingInfoRequest{
				Token:  msgs[0].Token,
				Staker: msgs[0].Staker,
			},
			response: &types.QueryStakingInfoResponse{StakingInfo: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryStakingInfoRequest{
				Token:  msgs[1].Token,
				Staker: msgs[1].Staker,
			},
			response: &types.QueryStakingInfoResponse{StakingInfo: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryStakingInfoRequest{
				Token:  strconv.Itoa(100000),
				Staker: strconv.Itoa(100000),
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
			response, err := keeper.StakingInfo(wctx, tc.request)
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

func TestStakingInfoQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNStakingInfo(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryStakingInfoAllRequest {
		return &types.QueryStakingInfoAllRequest{
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
			resp, err := keeper.StakingInfoAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StakingInfo), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.StakingInfo),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			//nolint:gosec // G115
			resp, err := keeper.StakingInfoAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StakingInfo), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.StakingInfo),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.StakingInfoAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		//nolint:gosec // G115
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.StakingInfo),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.StakingInfoAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
