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

func msgUpdateClass(creator string, id string) *types.MsgUpdateClass {
	return &types.MsgUpdateClass{
		Creator:     creator,
		Id:          id,
		Name:        sample.Name(),
		Symbol:      sample.Word(),
		Description: sample.Paragraph(),
		Uri:         sample.URL(),
		UriHash:     sample.Hash(),
		Data:        sample.JSON(),
	}
}

func mustUpdateClass(t testing.TB, ms types.MsgServer, ctx context.Context, ctrl *gomock.Controller, nftKeeper *testutil.MockNFTKeeper, updater string, classId string) {
	nftKeeper.EXPECT().UpdateClass(ctx, gomock.Any()).Return(nil)

	resp, err := ms.UpdateClass(ctx, msgUpdateClass(updater, classId))

	require.NoError(t, err)
	require.NotNil(t, resp)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	eventUpdateClass, err := utils.GetTypedEvent(sdkCtx, &types.EventUpdateClass{})

	require.NoError(t, err)
	require.NotNil(t, eventUpdateClass)
	require.Equal(t, classId, eventUpdateClass.Id)
}

func TestUpdateClass(t *testing.T) {
	ms, ctx, ctrl, nftKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	classId := mustCreateClass(t, ms, ctx, ctrl, nftKeeper, alice)

	mustUpdateClass(t, ms, ctx, ctrl, nftKeeper, alice, classId)
}

func TestUpdateClassNotFound(t *testing.T) {
	ms, ctx, ctrl, _ := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	resp, err := ms.UpdateClass(ctx, msgUpdateClass(alice, "1"))

	require.Nil(t, resp)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrNotFound)
}

func TestUpdateClassUnauthorized(t *testing.T) {
	ms, ctx, ctrl, nftKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	classId := mustCreateClass(t, ms, ctx, ctrl, nftKeeper, alice)

	resp, err := ms.UpdateClass(ctx, msgUpdateClass(bob, classId))

	require.Nil(t, resp)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrUnauthorized)
}
