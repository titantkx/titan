package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/x/farming/keeper"
	"github.com/titantkx/titan/x/farming/testutil"
	"github.com/titantkx/titan/x/farming/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context, *keeper.Keeper) {
	k, ctx := keepertest.FarmingKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx), k
}

func setupMsgServerWithMocks(t *testing.T) (types.MsgServer, context.Context, *gomock.Controller, *keeper.Keeper, *testutil.MockBankKeeper) {
	ctrl := gomock.NewController(t)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	k, ctx := keepertest.FarmingKeeperWithMocks(t, bankKeeper)
	msgSrv := keeper.NewMsgServerImpl(*k)
	return msgSrv, sdk.WrapSDKContext(ctx), ctrl, k, bankKeeper
}

func TestMsgServer(t *testing.T) {
	ms, ctx, k := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotNil(t, k)
}
