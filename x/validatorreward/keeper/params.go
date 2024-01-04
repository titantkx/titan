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
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}

	store.Set(types.ParamsKey, bz)

	return nil
}

// Rate returns the Rate param
func (k Keeper) Rate(ctx sdk.Context) sdk.Dec {
	params := k.GetParams(ctx)
	return params.Rate
}

// Operator returns the Operator param
func (k Keeper) Operator(ctx sdk.Context) sdk.AccAddress {
	params := k.GetParams(ctx)
	addr := sdk.MustAccAddressFromBech32(params.Operator)
	return addr
}
