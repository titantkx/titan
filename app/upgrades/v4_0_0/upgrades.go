package v4_0_0

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	pointertypes "github.com/titantkx/titan/x/pointer/types"
)

func CreateStoreUpgrade() *storetypes.StoreUpgrades {
	return &storetypes.StoreUpgrades{
		Added: []string{
			pointertypes.StoreKey,
		},
	}
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		//	add module package forward middleware
		//	NOTE: no need to do anything here.
		//	all process for add new module already done in `mm.RunMigrations`
		//	like call `InitGenesis` (with default genesis) instead of migration via `RegisterMigration`

		// Leave modules are as-is to avoid running InitGenesis.
		logger.Debug("running module migrations ...")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
