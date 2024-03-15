package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	sdkstakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/titantkx/titan/x/staking/types"
)

type Keeper struct {
	*stakingkeeper.Keeper
	distkeeper types.DistributionKeeper
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	ak sdkstakingtypes.AccountKeeper,
	bk sdkstakingtypes.BankKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		Keeper: stakingkeeper.NewKeeper(cdc, key, ak, bk, authority),
	}
}

func (k *Keeper) SetDistributionKeeper(dk types.DistributionKeeper) {
	k.distkeeper = dk
}
