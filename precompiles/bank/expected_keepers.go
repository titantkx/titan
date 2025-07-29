package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	pointertypes "github.com/titantkx/titan/x/pointer/types"
)

type PointerKeeper interface {
	GetErc20Native(
		ctx sdk.Context,
		tokenDenom string,
	) (val pointertypes.Erc20Native, found bool)
}
