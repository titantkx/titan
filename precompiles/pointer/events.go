package pointer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	EventPointerRegistered = "PointerRegistered"
)

func (p PrecompileExecutor) emitPointerRegisteredEvent(ctx sdk.Context,
	evm *vm.EVM, token string, contractAddr ethcommon.Address,
) error {
	event := p.ABI.Events[EventPointerRegistered] //nolint:staticcheck
	topics := []ethcommon.Hash{
		event.ID,
	}
	arguments := ethabi.Arguments{event.Inputs[0], event.Inputs[1], event.Inputs[2]}
	data, err := arguments.Pack(uint8(0), token, contractAddr)
	if err != nil {
		return err
	}

	evm.StateDB.AddLog(&ethtypes.Log{
		Address:     ethcommon.HexToAddress(precompileContractAddress),
		Topics:      topics,
		Data:        data,
		BlockNumber: uint64(ctx.BlockHeight()), //nolint:gosec
	})
	return nil
}
