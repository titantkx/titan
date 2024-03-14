package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/validatorreward/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RewardPool(goCtx context.Context, req *types.QueryRewardPoolRequest) (*types.QueryRewardPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pool := k.GetRewardPoolCoins(ctx)

	return &types.QueryRewardPoolResponse{Pool: pool}, nil
}
