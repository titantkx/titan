package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

func (k msgServer) SetOperator(goCtx context.Context, msg *types.MsgSetOperator) (*types.MsgSetOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate operator
	operatorAddr, err := sdk.AccAddressFromBech32(msg.Operator)
	if err != nil {
		return nil, err
	}

	if !operatorAddr.Equals(k.GetOperator(ctx)) {
		return nil, types.ErrForbidden.Wrapf("not allowed to set operator")
	}

	newOperatorAcc, err := sdk.AccAddressFromBech32(msg.NewOperator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrap("invalid new operator address")
	}

	k.Keeper.SetOperator(ctx, newOperatorAcc)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetOperator,
			sdk.NewAttribute(types.AttributeKeyOperator, msg.NewOperator),
		),
	)

	return &types.MsgSetOperatorResponse{}, nil
}
