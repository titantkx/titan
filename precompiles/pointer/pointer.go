package pointer

import (
	"embed"
	"errors"
	"fmt"
	"math"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	pcommon "github.com/titantkx/ethermint/precompiles/common"
)

const (
	PrecompileName   = "pointer"
	AddNativePointer = "addNativePointer"
)

const PointerAddress = "0x000000000000000000000000000000000000100b"

// Embed abi json file to the executable binary. Needed when importing as dependency.
//
//go:embed abi.json
var f embed.FS

var _ pcommon.PrecompileExecutor = &PrecompileExecutor{}

type PrecompileExecutor struct {
	bankKeeper BankKeeper

	AddNativePointerID []byte
}

func NewPrecompile(bankKeeper BankKeeper) *pcommon.Precompile {
	abi := pcommon.MustGetABI(f, "abi.json")

	p := &PrecompileExecutor{
		bankKeeper: bankKeeper,
	}

	for name, m := range abi.Methods {
		switch name {
		case AddNativePointer:
			p.AddNativePointerID = m.ID
		}
	}

	return pcommon.NewPrecompile(abi, ethcommon.HexToAddress(PointerAddress), p)
}

func (p *PrecompileExecutor) RequiredGas(input []byte, method *ethabi.Method) uint64 {
	return 0
}

func (p *PrecompileExecutor) Execute(
	ctx sdk.Context,
	stateDB vm.StateDB,
	method *ethabi.Method,
	caller common.Address,
	callingContract vm.ContractRef,
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
	case AddNativePointer:
		return p.AddNative(ctx, method, caller, args, value)
	default:
		return nil, fmt.Errorf("unknown method %s", method.Name)
	}
}

func (p PrecompileExecutor) AddNative(
	ctx sdk.Context,
	method *ethabi.Method,
	caller common.Address,
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
			symbol = denomUnit.Denom
			if len(denomUnit.Aliases) > 0 {
				name = denomUnit.Aliases[0]
			}
		}
	}

	// @todo call custom keeper to create pointer contract link to the native token
	_ = name
	_ = symbol

	// ret, err = method.Outputs.Pack(contractAddr)
	return ret, err
}
