package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/farming/types"
)

func (k msgServer) Harvest(goCtx context.Context, msg *types.MsgHarvest) (*types.MsgHarvestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	reward, found := k.GetReward(ctx, msg.Sender)
	if !found {
		return nil, types.ErrNoReward
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.MustAccAddressFromBech32(reward.Farmer), reward.Amount); err != nil {
		return nil, err
	}

	k.RemoveReward(ctx, msg.Sender)

	event := &types.EventHarvest{
		Sender: msg.Sender,
		Amount: reward.Amount,
	}

	if err := ctx.EventManager().EmitTypedEvent(event); err != nil {
		return nil, types.WrapInternalError(err)
	}

	return &types.MsgHarvestResponse{}, nil
}
