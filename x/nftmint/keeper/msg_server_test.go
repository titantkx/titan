package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/x/nftmint/keeper"
	"github.com/titantkx/titan/x/nftmint/testutil"
	"github.com/titantkx/titan/x/nftmint/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.NftmintKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func setupMsgServerWithMocks(t *testing.T) (types.MsgServer, context.Context, *gomock.Controller, *testutil.MockNFTKeeper) {
	ctrl := gomock.NewController(t)
	nftKeeper := testutil.NewMockNFTKeeper(ctrl)
	k, ctx := keepertest.NftmintKeeperWithMocks(t, nftKeeper)
	msgSrv := keeper.NewMsgServerImpl(*k)
	return msgSrv, sdk.WrapSDKContext(ctx), ctrl, nftKeeper
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}

func TestMsgServerWithMocks(t *testing.T) {
	ms, ctx, ctrl, nftKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotNil(t, ctrl)
	require.NotNil(t, nftKeeper)
}
