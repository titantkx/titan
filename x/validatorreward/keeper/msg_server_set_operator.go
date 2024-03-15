package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/titantkx/titan/x/validatorreward/types"
)

func (k msgServer) SetAuthority(goCtx context.Context, msg *types.MsgSetAuthority) (*types.MsgSetAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate authority
	authorityAddr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	if !authorityAddr.Equals(k.GetAuthority(ctx)) {
		return nil, types.ErrForbidden.Wrapf("not allowed to set authority")
	}

	newAuthorityAcc, err := sdk.AccAddressFromBech32(msg.NewAuthority)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrap("invalid new authority address")
	}

	k.Keeper.SetAuthority(ctx, newAuthorityAcc)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetAuthority,
			sdk.NewAttribute(types.AttributeKeyAuthority, msg.NewAuthority),
		),
	)

	return &types.MsgSetAuthorityResponse{}, nil
}
