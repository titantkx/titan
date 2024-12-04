package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/titantkx/titan/x/nftmint/types"
)

func (k msgServer) UpdateClass(goCtx context.Context, msg *types.MsgUpdateClass) (*types.MsgUpdateClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	mintingInfo, ok := k.GetMintingInfo(ctx, msg.Id)
	if !ok {
		return nil, types.WrapError(types.ErrNotFound, "class not found")
	}

	if msg.Creator != mintingInfo.Owner {
		return nil, types.WrapErrorf(types.ErrUnauthorized, "%s is not the owner of class %s", msg.Creator, msg.Id)
	}

	classData := types.MustNewAnyWithMetadata(msg.Data)

	class := nft.Class{
		Id:          msg.Id,
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		Description: msg.Description,
		Uri:         msg.Uri,
		UriHash:     msg.UriHash,
		Data:        classData,
	}

	if err := k.nftKeeper.UpdateClass(ctx, class); err != nil {
		return nil, types.WrapInternalError(err)
	}

	if err := ctx.EventManager().EmitTypedEvent(&types.EventUpdateClass{Id: msg.Id}); err != nil {
		return nil, types.WrapInternalError(err)
	}

	return &types.MsgUpdateClassResponse{}, nil
}
