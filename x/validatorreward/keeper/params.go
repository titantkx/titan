package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(clientCtx sdk.Context) (params types.Params) {
	store := clientCtx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}

	store.Set(types.ParamsKey, bz)

	return nil
}

// Rate returns the Rate param
func (k Keeper) GetRate(ctx sdk.Context) sdk.Dec {
	params := k.GetParams(ctx)
	return params.Rate
}

// SetRate sets the Rate param
func (k Keeper) SetRate(ctx sdk.Context, rate sdk.Dec) {
	params := k.GetParams(ctx)
	params.Rate = rate
	k.SetParams(ctx, params)
}

// Operator returns the Operator param
func (k Keeper) GetOperator(ctx sdk.Context) sdk.AccAddress {
	params := k.GetParams(ctx)
	addr := sdk.MustAccAddressFromBech32(params.Operator)
	return addr
}

// SetOperator sets the Operator param
func (k Keeper) SetOperator(ctx sdk.Context, operator sdk.AccAddress) {
	params := k.GetParams(ctx)
	params.Operator = operator.String()
	k.SetParams(ctx, params)
}
