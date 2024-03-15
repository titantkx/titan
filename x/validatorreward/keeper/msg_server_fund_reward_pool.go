package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/validatorreward/types"
)

func (k msgServer) FundRewardPool(goCtx context.Context, msg *types.MsgFundRewardPool) (*types.MsgFundRewardPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	depositer, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositer, types.ModuleName, msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgFundRewardPoolResponse{}, nil
}
