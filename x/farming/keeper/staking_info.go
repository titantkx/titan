package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/farming/types"
)

// SetStakingInfo set a specific stakingInfo in the store from its index
func (k Keeper) SetStakingInfo(ctx sdk.Context, stakingInfo types.StakingInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StakingInfoKeyPrefix))
	b := k.cdc.MustMarshal(&stakingInfo)
	store.Set(types.StakingInfoKey(
		stakingInfo.Token,
		stakingInfo.Staker,
	), b)
}

// GetStakingInfo returns a stakingInfo from its index
func (k Keeper) GetStakingInfo(
	ctx sdk.Context,
	token string,
	staker string,
) (val types.StakingInfo, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StakingInfoKeyPrefix))

	b := store.Get(types.StakingInfoKey(
		token,
		staker,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveStakingInfo removes a stakingInfo from the store
func (k Keeper) RemoveStakingInfo(
	ctx sdk.Context,
	token string,
	staker string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StakingInfoKeyPrefix))
	store.Delete(types.StakingInfoKey(
		token,
		staker,
	))
}

// GetAllStakingInfo returns all stakingInfo
func (k Keeper) GetAllStakingInfo(ctx sdk.Context, token string) (list []types.StakingInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StakingInfoKeyPrefix))

	prefix := []byte{}
	if token != "" {
		prefix = types.StakingInfoKey(token, "")
	}

	iterator := sdk.KVStorePrefixIterator(store, prefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.StakingInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
