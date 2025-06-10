package pointer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	pointertypes "github.com/titantkx/titan/x/pointer/types"
)

type BankKeeper interface {
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
}

type PointerKeeper interface { //nolint:revive
	DeployOrUpdateErc20NativePointer(
		ctx sdk.Context,
		evm *vm.EVM,
		token string, metadata pointertypes.ERCMetadata,
	) (contractAddr ethcommon.Address, err error)
}
