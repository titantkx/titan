package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/titantkx/titan/x/nftmint/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) MintingInfos(goCtx context.Context, req *types.QueryMintingInfosRequest) (*types.QueryMintingInfosResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var mintingInfos []types.MintingInfo
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	mintingInfoStore := prefix.NewStore(store, types.KeyPrefix(types.MintingInfoKeyPrefix))

	pageRes, err := query.Paginate(mintingInfoStore, req.Pagination, func(key []byte, value []byte) error {
		var mintingInfo types.MintingInfo
		if err := k.cdc.Unmarshal(value, &mintingInfo); err != nil {
			return err
		}

		mintingInfos = append(mintingInfos, mintingInfo)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMintingInfosResponse{MintingInfo: mintingInfos, Pagination: pageRes}, nil
}

func (k Keeper) MintingInfo(goCtx context.Context, req *types.QueryMintingInfoRequest) (*types.QueryMintingInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetMintingInfo(
		ctx,
		req.ClassId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryMintingInfoResponse{MintingInfo: val}, nil
}
