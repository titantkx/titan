package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/farming/types"
)

// SetFarm set a specific farm in the store from its index
func (k Keeper) SetFarm(ctx sdk.Context, farm types.Farm) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FarmKeyPrefix))
	b := k.cdc.MustMarshal(&farm)
	store.Set(types.FarmKey(
		farm.Token,
	), b)
}

// GetFarm returns a farm from its index
func (k Keeper) GetFarm(ctx sdk.Context, token string) (val types.Farm, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FarmKeyPrefix))

	b := store.Get(types.FarmKey(token))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveFarm removes a farm from the store
func (k Keeper) RemoveFarm(ctx sdk.Context, token string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FarmKeyPrefix))
	store.Delete(types.FarmKey(
		token,
	))
}

// GetAllFarm returns all farm
func (k Keeper) GetAllFarm(ctx sdk.Context) (list []types.Farm) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FarmKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Farm
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
