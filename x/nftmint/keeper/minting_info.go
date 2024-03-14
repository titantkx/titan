package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/nftmint/types"
)

// SetMintingInfo set a specific mintingInfo in the store from its index
func (k Keeper) SetMintingInfo(ctx sdk.Context, mintingInfo types.MintingInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MintingInfoKeyPrefix))
	b := k.cdc.MustMarshal(&mintingInfo)
	store.Set(types.MintingInfoKey(
		mintingInfo.ClassId,
	), b)
}

// GetMintingInfo returns a mintingInfo from its index
func (k Keeper) GetMintingInfo(
	ctx sdk.Context,
	classId string,
) (val types.MintingInfo, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MintingInfoKeyPrefix))

	b := store.Get(types.MintingInfoKey(
		classId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveMintingInfo removes a mintingInfo from the store
func (k Keeper) RemoveMintingInfo(
	ctx sdk.Context,
	classId string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MintingInfoKeyPrefix))
	store.Delete(types.MintingInfoKey(
		classId,
	))
}

// GetAllMintingInfo returns all mintingInfo
func (k Keeper) GetAllMintingInfo(ctx sdk.Context) (list []types.MintingInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MintingInfoKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MintingInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
