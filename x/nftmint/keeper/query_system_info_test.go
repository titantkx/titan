package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/nftmint/types"
)

func TestSystemInfoQuery(t *testing.T) {
	keeper, ctx := keepertest.NftmintKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestSystemInfo(keeper, ctx)
	tests := []struct {
		desc     string
		request  *types.QuerySystemInfoRequest
		response *types.QuerySystemInfoResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QuerySystemInfoRequest{},
			response: &types.QuerySystemInfoResponse{SystemInfo: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.SystemInfo(wctx, tc.request)
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
