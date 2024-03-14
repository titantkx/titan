package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/nftmint/testutil"
	"github.com/titantkx/titan/x/nftmint/types"
)

func msgCreateClass(creator string) *types.MsgCreateClass {
	return &types.MsgCreateClass{
		Creator:     creator,
		Name:        sample.Name(),
		Symbol:      sample.Word(),
		Description: sample.Paragraph(),
		Uri:         sample.URL(),
		UriHash:     sample.Hash(),
		Data:        sample.JSON(),
	}
}

func mustCreateClass(t testing.TB, ms types.MsgServer, ctx context.Context, ctrl *gomock.Controller, nftKeeper *testutil.MockNFTKeeper, creator string) string {
	nftKeeper.EXPECT().SaveClass(ctx, gomock.Any()).Return(nil)

	resp, err := ms.CreateClass(ctx, msgCreateClass(creator))

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Id)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	eventCreateClass, err := utils.GetTypedEvent(sdkCtx, &types.EventCreateClass{})

	require.NoError(t, err)
	require.NotNil(t, eventCreateClass)
	require.Equal(t, resp.Id, eventCreateClass.Id)
	require.Equal(t, creator, eventCreateClass.Owner)

	return resp.Id
}

func TestCreateClass(t *testing.T) {
	ms, ctx, ctrl, nftKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	mustCreateClass(t, ms, ctx, ctrl, nftKeeper, alice)
}
