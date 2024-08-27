package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func (server msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authorityMetadata, err := server.Keeper.GetAuthorityMetadata(ctx, msg.Amount.GetDenom())
	if err != nil {
		return nil, err
	}

	if msg.Sender != authorityMetadata.GetAdmin() {
		return nil, types.ErrUnauthorized
	}

	err = server.Keeper.burnFrom(ctx, msg.Amount, msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgBurn,
			sdk.NewAttribute(types.AttributeBurnFromAddress, msg.Sender),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
		),
	})

	return &types.MsgBurnResponse{}, nil
}
