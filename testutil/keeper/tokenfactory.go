package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/app"
	"github.com/titantkx/titan/x/tokenfactory/keeper"
	tokenfactoryutil "github.com/titantkx/titan/x/tokenfactory/testutil"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func TokenfactoryKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	return TokenfactoryKeeperWithMocks(t, nil, nil, nil, nil)
}

func TokenfactoryKeeperWithMocks(t testing.TB, accountKeeper *tokenfactoryutil.MockAccountKeeper, bankKeeper *tokenfactoryutil.MockBankKeeper, contractKeeper *tokenfactoryutil.MockContractKeeper, communityPoolKeeper *tokenfactoryutil.MockCommunityPoolKeeper) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		app.MaccPerms,
		accountKeeper,
		bankKeeper,
		contractKeeper,
		communityPoolKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	err := k.SetParams(ctx, types.DefaultParams())
	require.NoError(t, err)

	return k, ctx
}
