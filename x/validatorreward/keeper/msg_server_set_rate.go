package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/validatorreward/types"
)

func (k msgServer) SetRate(goCtx context.Context, msg *types.MsgSetRate) (*types.MsgSetRateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate authority
	authorityAddr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	if !authorityAddr.Equals(k.GetAuthority(ctx)) {
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
