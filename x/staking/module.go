package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	sdkstaking "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/titanlab/titan/x/staking/keeper"
)

// AppModuleBasic defines the basic application module used by the staking module.
type AppModule struct {
	sdkstaking.AppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec,
	keeper *keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	ls exported.Subspace,
) AppModule {
	return AppModule{
		AppModule: sdkstaking.NewAppModule(cdc, keeper.Keeper, ak, bk, ls),
	}
}
