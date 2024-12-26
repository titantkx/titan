package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/farming/types"
)

func (k msgServer) Unstake(goCtx context.Context, msg *types.MsgUnstake) (*types.MsgUnstakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, coin := range msg.Amount {
		stakingInfo, ok := k.GetStakingInfo(ctx, coin.Denom, msg.Sender)
		if !ok || stakingInfo.Amount.LT(coin.Amount) {
			return nil, types.ErrStakingBalanceNotEnough
		}

		stakingInfo.Amount = stakingInfo.Amount.Sub(coin.Amount)

		if stakingInfo.Amount.IsZero() {
			k.RemoveStakingInfo(ctx, stakingInfo.Token, stakingInfo.Staker)
		} else {
			k.SetStakingInfo(ctx, stakingInfo)
		}
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.MustAccAddressFromBech32(msg.Sender), msg.Amount); err != nil {
		return nil, err
	}

	event := &types.EventUnstake{
		Sender: msg.Sender,
		Amount: msg.Amount,
	}

	if err := ctx.EventManager().EmitTypedEvent(event); err != nil {
		return nil, types.WrapInternalError(err)
	}

	return &types.MsgUnstakeResponse{}, nil
}
