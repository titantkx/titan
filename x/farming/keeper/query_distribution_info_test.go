package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/farming/types"
)

func TestDistributionInfoQuery(t *testing.T) {
	keeper, ctx := keepertest.FarmingKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestDistributionInfo(keeper, ctx)
	tests := []struct {
		desc     string
		request  *types.QueryDistributionInfoRequest
		response *types.QueryDistributionInfoResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryDistributionInfoRequest{},
			response: &types.QueryDistributionInfoResponse{DistributionInfo: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.DistributionInfo(wctx, tc.request)
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
