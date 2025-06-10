package bank

import (
	"embed"
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	pcommon "github.com/titantkx/ethermint/precompiles/common"

	utils "github.com/titantkx/titan/utils"
)

const (
	PrecompileName = "bank"

	SendMethod        = "send"
	BalanceMethod     = "balance"
	AllBalancesMethod = "all_balances"
	NameMethod        = "name"
	SymbolMethod      = "symbol"
	DecimalsMethod    = "decimals"
	SupplyMethod      = "supply"
)

const PrecompileContractAddress = "0x0000000000000000000000000000000000001001"

// Embed abi json file to the executable binary. Needed when importing as dependency.
//
//go:embed abi.json
var f embed.FS

type CoinBalance struct {
	Amount *big.Int
	Denom  string
}

var _ pcommon.PrecompileExecutor = &PrecompileExecutor{}

type PrecompileExecutor struct {
	pointerKeeper PointerKeeper
	bankKeeper    bankkeeper.Keeper
}

func NewPrecompile(pointerKeeper PointerKeeper, bankKeeper bankkeeper.Keeper) *pcommon.Precompile {
	abi := pcommon.MustGetABI(f, "abi.json")

	p := &PrecompileExecutor{
		pointerKeeper: pointerKeeper,
		bankKeeper:    bankKeeper,
	}

	return pcommon.NewPrecompile(abi, ethcommon.HexToAddress(PrecompileContractAddress), p)
}

func (p *PrecompileExecutor) RequiredGas(input []byte, method *ethabi.Method) uint64 { //nolint:revive
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
	if isFromDelegateCall {
		return nil, errors.New("cannot call bank precompile from delegatecall")
	}

	switch method.Name {
	case SendMethod:
		return p.send(ctx, evm, method, caller, args, value, readOnly)
	case BalanceMethod:
		return p.balance(ctx, evm, method, caller, args, value)
	case AllBalancesMethod:
		return p.allBalances(ctx, evm, method, caller, args, value)
	case NameMethod:
		return p.name(ctx, evm, method, caller, args, value)
	case SymbolMethod:
		return p.symbol(ctx, evm, method, caller, args, value)
	case DecimalsMethod:
		return p.decimals(ctx, evm, method, caller, args, value)
	case SupplyMethod:
		return p.totalSupply(ctx, evm, method, caller, args, value)
	default:
		return nil, fmt.Errorf("unhandled method %s", method.Name)
	}
}

func (p PrecompileExecutor) send(
	ctx sdk.Context,
	_ *vm.EVM,
	method *ethabi.Method,
	caller ethcommon.Address,
	args []interface{},
	value *big.Int,
	readOnly bool,
) ([]byte, error) {
	if readOnly {
		return nil, errors.New("cannot call send from staticcall")
	}

	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, err
	}

	if err := pcommon.ValidateArgsLength(args, 4); err != nil {
		return nil, err
	}
	tokenDenom := args[2].(string)
	if tokenDenom == "" {
		return nil, fmt.Errorf("token denom is empty")
	}
	// this method not allow to send base token (atkx)
	if tokenDenom == utils.BaseDenom {
		return nil, fmt.Errorf("this method not allow to send base token")
	}

	erc20Native, found := p.pointerKeeper.GetErc20Native(ctx, tokenDenom)
	if !found {
		return nil, fmt.Errorf("erc20native %s not found", tokenDenom)
	}
	if erc20Native.Erc20Addr != caller.Hex() {
		return nil, fmt.Errorf("only pointer %s can call send %s", erc20Native.Erc20Addr, caller.Hex())
	}

	from, err := pcommon.GetCosmosAddressFromEVMAddressArg(args[0])
	if err != nil {
		return nil, err
	}
	to, err := pcommon.GetCosmosAddressFromEVMAddressArg(args[1])
	if err != nil {
		return nil, err
	}
	amount := args[3].(*big.Int)
	if amount.Cmp(big.NewInt(0)) == 0 {
		// short circuit
		bz, err := method.Outputs.Pack(true)
		return bz, err
	}

	msg := &banktypes.MsgSend{
		FromAddress: from.String(),
		ToAddress:   to.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(tokenDenom, sdk.NewIntFromBigInt(amount))),
	}

	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	msgSrv := bankkeeper.NewMsgServerImpl(p.bankKeeper)
	_, err = msgSrv.Send(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return nil, err
	}

	bz, err := method.Outputs.Pack(true)
	return bz, err
}

func (p PrecompileExecutor) balance(
	ctx sdk.Context,
	_ *vm.EVM,
	method *ethabi.Method,
	_ ethcommon.Address,
	args []interface{},
	value *big.Int,
) ([]byte, error) {
	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, err
	}

	if err := pcommon.ValidateArgsLength(args, 2); err != nil {
		return nil, err
	}

	addr, err := pcommon.GetCosmosAddressFromEVMAddressArg(args[0])
	if err != nil {
		return nil, err
	}
	denom := args[1].(string)
	if denom == "" {
		return nil, errors.New("invalid denom")
	}

	bz, err := method.Outputs.Pack(p.bankKeeper.GetBalance(ctx, addr, denom).Amount.BigInt())
	return bz, err
}

func (p PrecompileExecutor) allBalances(
	ctx sdk.Context,
	_ *vm.EVM,
	method *ethabi.Method,
	_ ethcommon.Address,
	args []interface{},
	value *big.Int,
) ([]byte, error) {
	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, err
	}

	if err := pcommon.ValidateArgsLength(args, 1); err != nil {
		return nil, err
	}

	addr, err := pcommon.GetCosmosAddressFromEVMAddressArg(args[0])
	if err != nil {
		return nil, err
	}

	coins := p.bankKeeper.GetAllBalances(ctx, addr)

	// convert to coin balance structs
	coinBalances := make([]CoinBalance, 0, len(coins))

	for _, coin := range coins {
		coinBalances = append(coinBalances, CoinBalance{
			Amount: coin.Amount.BigInt(),
			Denom:  coin.Denom,
		})
	}

	bz, err := method.Outputs.Pack(coinBalances)
	return bz, err
}

func (p PrecompileExecutor) name(
	ctx sdk.Context,
	_ *vm.EVM,
	method *ethabi.Method,
	_ ethcommon.Address,
	args []interface{},
	value *big.Int,
) ([]byte, error) {
	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, err
	}

	if err := pcommon.ValidateArgsLength(args, 1); err != nil {
		return nil, err
	}

	denom := args[0].(string)
	metadata, found := p.bankKeeper.GetDenomMetaData(ctx, denom)
	if !found {
		return nil, fmt.Errorf("denom %s not found", denom)
	}
	bz, err := method.Outputs.Pack(metadata.Name)
	return bz, err
}

func (p PrecompileExecutor) symbol(
	ctx sdk.Context,
	_ *vm.EVM,
	method *ethabi.Method,
	_ ethcommon.Address,
	args []interface{},
	value *big.Int,
) ([]byte, error) {
	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, err
	}

	if err := pcommon.ValidateArgsLength(args, 1); err != nil {
		return nil, err
	}

	denom := args[0].(string)
	metadata, found := p.bankKeeper.GetDenomMetaData(ctx, denom)
	if !found {
		return nil, fmt.Errorf("denom %s not found", denom)
	}
	bz, err := method.Outputs.Pack(metadata.Symbol)
	return bz, err
}

func (p PrecompileExecutor) decimals(
	_ sdk.Context,
	_ *vm.EVM,
	method *ethabi.Method,
	_ ethcommon.Address,
	_ []interface{},
	value *big.Int,
) ([]byte, error) {
	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, err
	}

	// all native tokens are integer-based, returns decimals 0
	bz, err := method.Outputs.Pack(uint8(0))
	return bz, err
}

func (p PrecompileExecutor) totalSupply(
	ctx sdk.Context,
	_ *vm.EVM,
	method *ethabi.Method,
	_ ethcommon.Address,
	args []interface{},
	value *big.Int,
) ([]byte, error) {
	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, err
	}

	if err := pcommon.ValidateArgsLength(args, 1); err != nil {
		return nil, err
	}

	denom := args[0].(string)
	coin := p.bankKeeper.GetSupply(ctx, denom)
	bz, err := method.Outputs.Pack(coin.Amount.BigInt())
	return bz, err
}
