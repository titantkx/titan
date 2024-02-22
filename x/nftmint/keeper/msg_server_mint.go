package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/tokenize-titan/titan/x/nftmint/types"
)

func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	mintingInfo, ok := k.GetMintingInfo(ctx, msg.ClassId)
	if !ok {
		return nil, types.WrapError(types.ErrNotFound, "class not found")
	}

	if msg.Creator != mintingInfo.Owner {
		return nil, types.WrapErrorf(types.ErrUnauthorized, "%s is not the owner of class %s", msg.Creator, msg.ClassId)
	}

	tokenId := strconv.FormatUint(mintingInfo.NextTokenId, 10)
	tokenData := types.MustNewAnyWithMetadata(msg.Data)
	receiver := sdk.MustAccAddressFromBech32(msg.Receiver)

	token := nft.NFT{
		ClassId: msg.ClassId,
		Id:      tokenId,
		Uri:     msg.Uri,
		UriHash: msg.UriHash,
		Data:    tokenData,
	}

	if err := k.nftKeeper.Mint(ctx, token, receiver); err != nil {
		return nil, types.WrapInternalError(err)
	}

	mintingInfo.NextTokenId++
	k.SetMintingInfo(ctx, mintingInfo)

	return &types.MsgMintResponse{Id: tokenId}, nil
}
