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
	"github.com/titantkx/titan/x/nftmint/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestMintingInfoQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.NftmintKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNMintingInfo(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryMintingInfoRequest
		response *types.QueryMintingInfoResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryMintingInfoRequest{
				ClassId: msgs[0].ClassId,
			},
			response: &types.QueryMintingInfoResponse{MintingInfo: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryMintingInfoRequest{
				ClassId: msgs[1].ClassId,
			},
			response: &types.QueryMintingInfoResponse{MintingInfo: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryMintingInfoRequest{
				ClassId: strconv.Itoa(100000),
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
			response, err := keeper.MintingInfo(wctx, tc.request)
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

func TestMintingInfoQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.NftmintKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNMintingInfo(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryMintingInfosRequest {
		return &types.QueryMintingInfosRequest{
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
			resp, err := keeper.MintingInfos(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.MintingInfo), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.MintingInfo),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			//nolint:gosec // G115
			resp, err := keeper.MintingInfos(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.MintingInfo), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.MintingInfo),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.MintingInfos(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, uint64(len(msgs)), resp.Pagination.Total)
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.MintingInfo),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.MintingInfos(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
