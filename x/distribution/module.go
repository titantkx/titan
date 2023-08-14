package distribution

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	sdkdistribution "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/titanlab/titan/x/distribution/keeper"
)

type AppModule struct {
	sdkdistribution.AppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec, keeper keeper.Keeper, accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper, ss exported.Subspace,
) AppModule {
	return AppModule{
		AppModule: sdkdistribution.NewAppModule(cdc, keeper.Keeper, accountKeeper, bankKeeper, stakingKeeper, ss),
	}
}
