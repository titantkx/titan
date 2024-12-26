package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/titantkx/titan/x/farming/types"
)

// AddReward add reward for a farmer
func (k Keeper) AddReward(ctx sdk.Context, farmer string, amount sdk.Coins) {
	reward, found := k.GetReward(ctx, farmer)
	if !found {
		reward = types.Reward{
			Farmer: farmer,
			Amount: sdk.NewCoins(),
		}
	}

	reward.Amount = reward.Amount.Add(amount...)
	k.SetReward(ctx, reward)
}

// SetReward set a specific reward in the store from its index
func (k Keeper) SetReward(ctx sdk.Context, reward types.Reward) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RewardKeyPrefix))
	b := k.cdc.MustMarshal(&reward)
	store.Set(types.RewardKey(
		reward.Farmer,
	), b)
}

// GetReward returns a reward from its index
func (k Keeper) GetReward(ctx sdk.Context, farmer string) (val types.Reward, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RewardKeyPrefix))

	b := store.Get(types.RewardKey(
		farmer,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveReward removes a reward from the store
func (k Keeper) RemoveReward(ctx sdk.Context, farmer string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RewardKeyPrefix))
	store.Delete(types.RewardKey(
		farmer,
	))
}

// GetAllReward returns all reward
func (k Keeper) GetAllReward(ctx sdk.Context) (list []types.Reward) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RewardKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Reward
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
