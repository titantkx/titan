package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/titantkx/titan/x/pointer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Erc20NativeAll(goCtx context.Context, req *types.QueryErc20NativeAllRequest) (*types.QueryErc20NativeAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var erc20Natives []types.Erc20Native
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	erc20NativeStore := prefix.NewStore(store, types.KeyPrefix(types.Erc20NativeKeyPrefix))

	pageRes, err := query.Paginate(erc20NativeStore, req.Pagination, func(_ []byte, value []byte) error {
		var erc20Native types.Erc20Native
		if err := k.cdc.Unmarshal(value, &erc20Native); err != nil {
			return err
		}

		erc20Natives = append(erc20Natives, erc20Native)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryErc20NativeAllResponse{Erc20Native: erc20Natives, Pagination: pageRes}, nil
}

func (k Keeper) Erc20Native(goCtx context.Context, req *types.QueryErc20NativeRequest) (*types.QueryErc20NativeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetErc20Native(
		ctx,
		req.TokenDenom,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryErc20NativeResponse{Erc20Native: val}, nil
}
