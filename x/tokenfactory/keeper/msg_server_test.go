package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/x/tokenfactory/keeper"
	"github.com/titantkx/titan/x/tokenfactory/testutil"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.TokenfactoryKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func setupMsgServerWithMocks(t *testing.T) (types.MsgServer, context.Context, *gomock.Controller, *testutil.MockAccountKeeper, *testutil.MockBankKeeper, *testutil.MockContractKeeper, *testutil.MockCommunityPoolKeeper) {
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	contractKeeper := testutil.NewMockContractKeeper(ctrl)
	communityPoolKeeper := testutil.NewMockCommunityPoolKeeper(ctrl)
	k, ctx := keepertest.TokenfactoryKeeperWithMocks(t, accountKeeper, bankKeeper, contractKeeper, communityPoolKeeper)
	msgSrv := keeper.NewMsgServerImpl(*k)
	return msgSrv, sdk.WrapSDKContext(ctx), ctrl, accountKeeper, bankKeeper, contractKeeper, communityPoolKeeper
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}

func TestMsgServerWithMocks(t *testing.T) {
	ms, ctx, ctrl, accountKeeper, bankKeeper, contractKeeper, communityPoolKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotNil(t, ctrl)
	require.NotNil(t, accountKeeper)
	require.NotNil(t, bankKeeper)
	require.NotNil(t, contractKeeper)
	require.NotNil(t, communityPoolKeeper)
}
