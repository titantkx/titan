package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/nftmint/types"
)

func (k msgServer) TransferClass(goCtx context.Context, msg *types.MsgTransferClass) (*types.MsgTransferClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	mintingInfo, ok := k.GetMintingInfo(ctx, msg.ClassId)
	if !ok {
		return nil, types.WrapError(types.ErrNotFound, "class not found")
	}

	if msg.Creator != mintingInfo.Owner {
		return nil, types.WrapErrorf(types.ErrUnauthorized, "%s is not the owner of class %s", msg.Creator, msg.ClassId)
	}

	mintingInfo.Owner = msg.Receiver
	k.SetMintingInfo(ctx, mintingInfo)

	ctx.EventManager().EmitTypedEvent(&types.EventTransferClass{
		Id:       msg.ClassId,
		OldOwner: msg.Creator,
		NewOwner: msg.Receiver,
	})

	return &types.MsgTransferClassResponse{}, nil
}
