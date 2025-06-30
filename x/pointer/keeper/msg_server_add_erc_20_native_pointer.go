package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/titantkx/titan/x/pointer/types"
)

func (k msgServer) AddERC20NativePointer(goCtx context.Context, msg *types.MsgAddERC20NativePointer) (*types.MsgAddERC20NativePointerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate the authority
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}
	if authority.String() != k.authority {
		return nil, types.WrapError(sdkerrors.ErrUnauthorized, "only gov is allowed to perform this operation")
	}

	var contractAddr ethcommon.Address
	err = k.evmKeeper.RunWithOneOffEVMInstance(ctx, ethcommon.BytesToAddress(authority),
		func(evm *vm.EVM) error {
			contractAddr, err = k.DeployOrUpdateErc20NativePointer(ctx, evm, msg.Token, types.ERCMetadata{
				Name: msg.Name, Symbol: msg.Symbol, Decimals: uint8(msg.Decimals), //nolint:gosec
			})
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgAddERC20NativePointerResponse{
		ContractAddress: contractAddr.Hex(),
	}, nil
}
