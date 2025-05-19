package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/pointer/types"
)

func (k msgServer) AddERC20NativePointer(goCtx context.Context, msg *types.MsgAddERC20NativePointer) (*types.MsgAddERC20NativePointerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgAddERC20NativePointerResponse{}, nil
}
