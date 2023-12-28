package posthandler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
)

func NewPostHandler(
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	feegrantKeeper feegrantkeeper.Keeper,
) sdk.PostHandler {
	postDecorators := []sdk.PostDecorator{
		// The refund gas remaining decorator must be the last decorator in this
		// list.
		NewRefundGasRemainingDecorator(accountKeeper, bankKeeper, feegrantKeeper),
	}

	return sdk.ChainPostDecorators(postDecorators...)
}
