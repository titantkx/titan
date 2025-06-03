package pointer

import (
	"embed"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	pcommon "github.com/titantkx/ethermint/precompiles/common"

	pointertypes "github.com/titantkx/titan/x/pointer/types"
)

const (
	PrecompileName = "pointer"

	AddNativePointerMethod = "addNativePointer"
)

const precompileContractAddress = "0x000000000000000000000000000000000000100b"

// Embed abi json file to the executable binary. Needed when importing as dependency.
//
//go:embed abi.json
var f embed.FS

var _ pcommon.PrecompileExecutor = &PrecompileExecutor{}

type PrecompileExecutor struct {
	ethabi.ABI

	pointerKeeper PointerKeeper
	bankKeeper    BankKeeper
}

func NewPrecompile(pointerKeeper PointerKeeper, bankKeeper BankKeeper) *pcommon.Precompile {
	abi := pcommon.MustGetABI(f, "abi.json")

	p := &PrecompileExecutor{
		ABI:           abi,
		pointerKeeper: pointerKeeper,
		bankKeeper:    bankKeeper,
	}

	return pcommon.NewPrecompile(abi, ethcommon.HexToAddress(precompileContractAddress), p)
}

func (p *PrecompileExecutor) RequiredGas(input []byte, method *ethabi.Method) uint64 {
	return 0
}

func (p *PrecompileExecutor) Execute(
	ctx sdk.Context,
	evm *vm.EVM,
	method *ethabi.Method,
	caller ethcommon.Address,
	callingContract vm.ContractRef, //nolint:revive
	args []interface{},
	value *big.Int,
	readOnly bool,
	isFromDelegateCall bool,
) (ret []byte, err error) {
	if readOnly {
		return nil, errors.New("cannot call pointer precompile from staticcall")
	}
	if isFromDelegateCall {
		return nil, errors.New("cannot call pointer precompile from delegatecall")
	}

	switch method.Name {
	case AddNativePointerMethod:
		return p.AddNative(ctx, evm, method, caller, args, value)
	default:
		return nil, fmt.Errorf("unknown method %s", method.Name)
	}
}

func (p PrecompileExecutor) AddNative(
	ctx sdk.Context,
	evm *vm.EVM,
	method *ethabi.Method,
	caller ethcommon.Address, //nolint:revive
	args []interface{},
	value *big.Int,
) (ret []byte, err error) {
	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, err
	}
	if err := pcommon.ValidateArgsLength(args, 1); err != nil {
		return nil, err
	}

	token := args[0].(string)

	metadata, metadataExists := p.bankKeeper.GetDenomMetaData(ctx, token)
	if !metadataExists {
		return nil, fmt.Errorf("denom %s does not have metadata stored", token)
	}
	name := metadata.Name
	symbol := metadata.Symbol
	var decimals uint8
	for _, denomUnit := range metadata.DenomUnits {
		if denomUnit.Exponent > uint32(decimals) && denomUnit.Exponent <= math.MaxUint8 {
			decimals = uint8(denomUnit.Exponent)
			name = denomUnit.Denom
			symbol = strings.ToUpper(denomUnit.Denom)
			if len(denomUnit.Aliases) > 0 {
				name = denomUnit.Aliases[0]
			}
		}
	}

	// call pointer keeper to create pointer contract link to the native token
	contractAddr, err := p.pointerKeeper.DeployOrUpdateErc20NativePointer(ctx, evm, token, pointertypes.ERCMetadata{
		Name: name, Symbol: symbol, Decimals: decimals,
	})
	if err != nil {
		return nil, err
	}

	// emit event
	err = p.emitPointerRegisteredEvent(ctx, evm, token, contractAddr)
	if err != nil {
		return nil, err
	}

	ret, err = method.Outputs.Pack(contractAddr)
	return ret, err
}
