package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tokenize-titan/titan/utils"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

func (k Keeper) DistributeTokens(ctx sdk.Context) {
	// get current balance of reward pool
	validatorRewardAccount := k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
	currentBalance := k.bankKeeper.GetBalance(ctx, validatorRewardAccount.GetAddress(), utils.BondDenom)

	//@todo
	fmt.Println("currentBalance", currentBalance)
}
