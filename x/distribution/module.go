package distribution

import (
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkdistribution "github.com/cosmos/cosmos-sdk/x/distribution"
	sdkdistributionexported "github.com/cosmos/cosmos-sdk/x/distribution/exported"
	sdkdistributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/tokenize-titan/titan/x/distribution/keeper"
)

// AppModuleBasic defines the basic application module used by the distribution module.
type AppModule struct {
	sdkdistribution.AppModule

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec, keeper keeper.Keeper, accountKeeper sdkdistributiontypes.AccountKeeper,
	bankKeeper sdkdistributiontypes.BankKeeper, stakingKeeper sdkdistributiontypes.StakingKeeper, ss sdkdistributionexported.Subspace,
) AppModule {
	return AppModule{
		AppModule: sdkdistribution.NewAppModule(cdc, keeper.Keeper, accountKeeper, bankKeeper, stakingKeeper, ss),

		keeper: keeper,
	}
}

// BeginBlock returns the begin blocker for the distribution module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}
