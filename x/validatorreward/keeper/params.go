package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.Rate(ctx),
		k.Operator(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// Rate returns the Rate param
func (k Keeper) Rate(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyRate, &res)
	return
}

// Operator returns the Operator param
func (k Keeper) Operator(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyOperator, &res)
	return
}
