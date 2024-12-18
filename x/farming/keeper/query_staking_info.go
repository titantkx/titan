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

func (k Keeper) StakingInfoAll(goCtx context.Context, req *types.QueryStakingInfoAllRequest) (*types.QueryStakingInfoAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var stakingInfos []types.StakingInfo
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	stakingInfoStore := prefix.NewStore(store, types.KeyPrefix(types.StakingInfoKeyPrefix))
	if req.Token != "" {
		stakingInfoStore = prefix.NewStore(stakingInfoStore, types.StakingInfoKey(req.Token, ""))
	}

	pageRes, err := query.Paginate(stakingInfoStore, req.Pagination, func(_ []byte, value []byte) error {
		var stakingInfo types.StakingInfo
		if err := k.cdc.Unmarshal(value, &stakingInfo); err != nil {
			return err
		}

		stakingInfos = append(stakingInfos, stakingInfo)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryStakingInfoAllResponse{StakingInfo: stakingInfos, Pagination: pageRes}, nil
}

func (k Keeper) StakingInfo(goCtx context.Context, req *types.QueryStakingInfoRequest) (*types.QueryStakingInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetStakingInfo(
		ctx,
		req.Token,
		req.Staker,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryStakingInfoResponse{StakingInfo: val}, nil
}
