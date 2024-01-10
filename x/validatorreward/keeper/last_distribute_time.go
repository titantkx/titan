package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

// SetLastDistributeTime set the last distribute time
func (k Keeper) SetLastDistributeTime(ctx sdk.Context, value time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&types.TimestampProto{Timestamp: &value})
	store.Set(types.KeyPrefix(types.LastDistributeTimeKey), bz)
}

func (k Keeper) GetLastDistributeTime(ctx sdk.Context) (value time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefix(types.LastDistributeTimeKey))
	if bz == nil {
		return time.Time{}
	}
	var timestamp types.TimestampProto
	k.cdc.MustUnmarshal(bz, &timestamp)
	return *timestamp.Timestamp
}

func (k Keeper) RemoveLastDistributeTime(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyPrefix(types.LastDistributeTimeKey))
}
