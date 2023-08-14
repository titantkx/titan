package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DistributionKeeper expected distribution keeper (noalias)
type DistributionKeeper interface {
	GetFeePoolCommunityCoins(ctx sdk.Context) sdk.DecCoins
	GetValidatorOutstandingRewardsCoins(ctx sdk.Context, val sdk.ValAddress) sdk.DecCoins
	FundCommunityPoolFromModule(ctx sdk.Context, amount sdk.Coins, senderModule string) error
}
