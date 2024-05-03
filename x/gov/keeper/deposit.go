package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// DeleteAndContributeDeposits deletes and contribute all the deposits on a specific proposal.
func (k Keeper) DeleteAndContributeDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)

	k.IterateDeposits(ctx, proposalID, func(deposit v1.Deposit) bool {
		err := k.distKeeper.FundCommunityPoolFromModule(ctx, deposit.Amount, types.ModuleName)
		if err != nil {
			panic(err)
		}

		depositor := sdk.MustAccAddressFromBech32(deposit.Depositor)

		store.Delete(types.DepositKey(proposalID, depositor))
		return false
	})
}
