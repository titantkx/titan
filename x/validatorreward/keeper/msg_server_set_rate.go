package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

func (k msgServer) SetRate(goCtx context.Context, msg *types.MsgSetRate) (*types.MsgSetRateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate operator
	operatorAddr, err := sdk.AccAddressFromBech32(msg.Operator)
	if err != nil {
		return nil, err
	}

	if !operatorAddr.Equals(k.GetOperator(ctx)) {
		return nil, types.ErrForbidden.Wrapf("not allowed to set rate")
	}

	k.Keeper.SetRate(ctx, msg.Rate)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetRate,
			sdk.NewAttribute(types.AttributeKeyRate, msg.Rate.String()),
		),
	)

	return &types.MsgSetRateResponse{}, nil
}
