//nolint:dupl
package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/titantkx/titan/x/farming/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RewardAll(goCtx context.Context, req *types.QueryRewardAllRequest) (*types.QueryRewardAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var rewards []types.Reward
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	rewardStore := prefix.NewStore(store, types.KeyPrefix(types.RewardKeyPrefix))

	pageRes, err := query.Paginate(rewardStore, req.Pagination, func(_ []byte, value []byte) error {
		var reward types.Reward
		if err := k.cdc.Unmarshal(value, &reward); err != nil {
			return err
		}

		rewards = append(rewards, reward)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRewardAllResponse{Reward: rewards, Pagination: pageRes}, nil
}

func (k Keeper) Reward(goCtx context.Context, req *types.QueryRewardRequest) (*types.QueryRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetReward(
		ctx,
		req.Farmer,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryRewardResponse{Reward: val}, nil
}
