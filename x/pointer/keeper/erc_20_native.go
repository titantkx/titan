package keeper

import (
	"fmt"
	"math/big"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	etherminttypes "github.com/titantkx/ethermint/types"

	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/pointer/artifacts"
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
		return sdkerrors.Wrap(errortypes.ErrInvalidAddress, "invalid zero erc20 address")
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

func (k Keeper) DeployOrUpdateErc20NativePointer(
	ctx sdk.Context,
	evm *vm.EVM,
	token string, metadata types.ERCMetadata,
) (contractAddr ethcommon.Address, err error) {
	pointerModuleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	pointerModuleEthAddr := ethcommon.BytesToAddress(pointerModuleAddr)
	pointerType := "erc20native"

	// not allow to create pointer for base token (atkx)
	if token == utils.BaseDenom {
		return ethcommon.Address{}, fmt.Errorf("cannot create pointer for base token %s", utils.BaseDenom)
	}

	var bin []byte
	args := []interface{}{
		token, metadata.Name, metadata.Symbol, metadata.Decimals,
	}
	bin, err = artifacts.GetParsedABI(pointerType).Pack("", args...)
	if err != nil {
		panic(err)
	}
	bin = append(artifacts.GetBin(pointerType), bin...)

	erc20Native, found := k.GetErc20Native(ctx, token)

	suppliedGas := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumedToLimit()
	var remainingGas uint64

	if found {
		var ret []byte
		// pointer contract already exited, update it code
		contractAddr, err = types.GetErc20AddressFromString(erc20Native.Erc20Addr)
		if err != nil {
			return
		}
		ret, remainingGas, err = evm.GetDeploymentCode(vm.AccountRef(pointerModuleEthAddr), bin, suppliedGas, big.NewInt(0), contractAddr)
		evm.StateDB.SetCode(contractAddr, ret)
	} else {
		// deploy new pointer contract
		// @todo maybe need to override set nonce of EVM like ethermint did
		_, contractAddr, remainingGas, err = evm.Create(vm.AccountRef(pointerModuleEthAddr), bin, suppliedGas, big.NewInt(0))
	}
	if err != nil {
		return
	}

	ctx.GasMeter().ConsumeGas(suppliedGas-remainingGas, "erc20native contract deploy or update")
	// set erc20native contract address
	k.SetErc20Native(ctx, types.Erc20Native{
		TokenDenom: token,
		Erc20Addr:  contractAddr.Hex(),
	})
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypePointerRegistered, sdk.NewAttribute(types.AttributeKeyPointerType, pointerType),
		sdk.NewAttribute(types.AttributeKeyPointerAddress, contractAddr.Hex()), sdk.NewAttribute(types.AttributeKeyPointee, token)))
	return
}
