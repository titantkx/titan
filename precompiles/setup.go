package precompiles

import (
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/titantkx/ethermint/x/evm/vm"
	"github.com/titantkx/titan/precompiles/pointer"
	pointerkeeper "github.com/titantkx/titan/x/pointer/keeper"
)

func GetCustomPrecompiles(
	pointerKeeper pointerkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) vm.PrecompiledContracts {
	pointerP := pointer.NewPrecompile(pointerKeeper, bankKeeper)

	return vm.PrecompiledContracts{
		pointerP.Address(): pointerP,
	}
}
