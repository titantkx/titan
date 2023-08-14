package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkdistributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type Keeper struct {
	sdkdistributionkeeper.Keeper
	bankKeeper types.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey,
	ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper,
	feeCollectorName string, authority string,
) Keeper {
	return Keeper{
		Keeper:     sdkdistributionkeeper.NewKeeper(cdc, key, ak, bk, sk, feeCollectorName, authority),
		bankKeeper: bk,
	}
}

// FundCommunityPool allows an module to directly fund the community fund pool.
// The amount is first added to the distribution module account and then directly
// added to the pool. An error is returned if the amount cannot be sent to the
// module account.
func (k Keeper) FundCommunityPoolFromModule(ctx sdk.Context, amount sdk.Coins, senderModule string) error {
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, senderModule, types.ModuleName, amount); err != nil {
		return err
	}

	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(amount...)...)
	k.SetFeePool(ctx, feePool)

	return nil
}
