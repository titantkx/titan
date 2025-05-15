package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	etherminttypes "github.com/titantkx/ethermint/types"

	"github.com/titantkx/titan/x/pointer/types"
)

// SetErc20Native set a specific erc20Native in the store from its index
func (k Keeper) SetErc20Native(ctx sdk.Context, erc20Native types.Erc20Native) error {
	//  validate `erc20Native.Erc20Addr` before storing
	erc20Addr, err := types.GetErc20AddressFromString(erc20Native.Erc20Addr)
	if err != nil {
		return err
	}

	erc20Native.Erc20Addr = erc20Addr.Hex()

	if etherminttypes.IsZeroAddress(erc20Native.Erc20Addr) {
		return sdkerrors.Wrap(errortypes.ErrInvalidAddress, "invalid checksum for erc20 address")
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Erc20NativeKeyPrefix))
	b := k.cdc.MustMarshal(&erc20Native)
	store.Set(types.Erc20NativeKey(
		erc20Native.TokenDenom,
	), b)

	return nil
}

// GetErc20Native returns a erc20Native from its index
func (k Keeper) GetErc20Native(
	ctx sdk.Context,
	tokenDenom string,
) (val types.Erc20Native, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Erc20NativeKeyPrefix))

	b := store.Get(types.Erc20NativeKey(
		tokenDenom,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveErc20Native removes a erc20Native from the store
func (k Keeper) RemoveErc20Native(
	ctx sdk.Context,
	tokenDenom string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Erc20NativeKeyPrefix))
	store.Delete(types.Erc20NativeKey(
		tokenDenom,
	))
}

// GetAllErc20Native returns all erc20Native
func (k Keeper) GetAllErc20Native(ctx sdk.Context) (list []types.Erc20Native) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Erc20NativeKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Erc20Native
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
