package validatorreward

import (
	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tokenize-titan/titan/x/validatorreward/keeper"
)

// BeginBlocker add validator reward to `feeCollector`
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// @todo
	k.DistributeTokens(ctx)
}
