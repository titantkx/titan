package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/farming/types"
)

func (k msgServer) AddReward(goCtx context.Context, msg *types.MsgAddReward) (*types.MsgAddRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)

	if params.AddRewardGas > 0 {
		ctx.GasMeter().ConsumeGas(params.AddRewardGas, "add reward")
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(msg.Sender), types.ModuleName, msg.Amount); err != nil {
		return nil, err
	}

	farm, ok := k.GetFarm(ctx, msg.Token)
	if !ok {
		farm = types.Farm{Token: msg.Token}
	}

	if msg.StartTime.IsZero() {
		msg.StartTime = ctx.BlockTime()
	}

	if msg.StartTime.Before(ctx.BlockTime()) {
		return nil, types.WrapErrorf(types.ErrInvalidTime, "reward must start after %s", ctx.BlockTime().Format(time.RFC3339))
	}

	if !msg.EndTime.After(msg.StartTime) {
		return nil, types.WrapErrorf(types.ErrInvalidTime, "reward must end after %s", msg.StartTime.Format(time.RFC3339))
	}

	farm.Rewards = append(farm.Rewards, &types.FarmReward{
		Sender:    msg.Sender,
		Amount:    msg.Amount,
		EndTime:   msg.EndTime,
		StartTime: msg.StartTime,
	})

	k.SetFarm(ctx, farm)

	event := &types.EventAddReward{
		Sender:    msg.Sender,
		Token:     msg.Token,
		Amount:    msg.Amount,
		EndTime:   msg.EndTime,
		StartTime: msg.StartTime,
	}

	if err := ctx.EventManager().EmitTypedEvent(event); err != nil {
		return nil, types.WrapInternalError(err)
	}

	return &types.MsgAddRewardResponse{}, nil
}
