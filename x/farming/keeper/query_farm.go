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

func (k Keeper) FarmAll(goCtx context.Context, req *types.QueryFarmAllRequest) (*types.QueryFarmAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var farms []types.Farm
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	farmStore := prefix.NewStore(store, types.KeyPrefix(types.FarmKeyPrefix))

	pageRes, err := query.Paginate(farmStore, req.Pagination, func(_ []byte, value []byte) error {
		var farm types.Farm
		if err := k.cdc.Unmarshal(value, &farm); err != nil {
			return err
		}

		farms = append(farms, farm)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryFarmAllResponse{Farm: farms, Pagination: pageRes}, nil
}

func (k Keeper) Farm(goCtx context.Context, req *types.QueryFarmRequest) (*types.QueryFarmResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetFarm(
		ctx,
		req.Token,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryFarmResponse{Farm: val}, nil
}
