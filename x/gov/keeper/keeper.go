package keeper

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdkgovkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	sdkgovtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/tokenize-titan/titan/x/gov/types"
)

type Keeper struct {
	*sdkgovkeeper.Keeper

	bankKeeper sdkgovtypes.BankKeeper
	distKeeper types.DistributionKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey storetypes.StoreKey
}

func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, authKeeper sdkgovtypes.AccountKeeper,
	bankKeeper sdkgovtypes.BankKeeper, sk sdkgovtypes.StakingKeeper, distKeeper types.DistributionKeeper,
	router *baseapp.MsgServiceRouter, config sdkgovtypes.Config, authority string,
) *Keeper {
	return &Keeper{
		Keeper:     sdkgovkeeper.NewKeeper(cdc, key, authKeeper, bankKeeper, sk, router, config, authority),
		bankKeeper: bankKeeper,
		distKeeper: distKeeper,
		storeKey:   key,
	}
}
