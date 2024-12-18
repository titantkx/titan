package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/farming/types"
)

// SetDistributionInfo set distributionInfo in the store
func (k Keeper) SetDistributionInfo(ctx sdk.Context, distributionInfo types.DistributionInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DistributionInfoKey))
	b := k.cdc.MustMarshal(&distributionInfo)
	store.Set([]byte{0}, b)
}

// GetDistributionInfo returns distributionInfo
func (k Keeper) GetDistributionInfo(ctx sdk.Context) (val types.DistributionInfo, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DistributionInfoKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
