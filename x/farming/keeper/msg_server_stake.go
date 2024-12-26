package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/farming/types"
)

func (k msgServer) Stake(goCtx context.Context, msg *types.MsgStake) (*types.MsgStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(msg.Sender), types.ModuleName, msg.Amount); err != nil {
		return nil, err
	}

	for _, coin := range msg.Amount {
		stakingInfo, ok := k.GetStakingInfo(ctx, coin.Denom, msg.Sender)
		if !ok {
			stakingInfo = types.StakingInfo{
				Token:  coin.Denom,
				Staker: msg.Sender,
				Amount: math.NewInt(0),
			}
		}
		stakingInfo.Amount = stakingInfo.Amount.Add(coin.Amount)
		k.SetStakingInfo(ctx, stakingInfo)
	}

	event := &types.EventStake{
		Sender: msg.Sender,
		Amount: msg.Amount,
	}

	if err := ctx.EventManager().EmitTypedEvent(event); err != nil {
		return nil, types.WrapInternalError(err)
	}

	return &types.MsgStakeResponse{}, nil
}
