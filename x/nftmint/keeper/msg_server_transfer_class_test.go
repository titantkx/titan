package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/nftmint/types"
)

//nolint:revive	// ctx at third position .
func mustTransferClass(t testing.TB, ms types.MsgServer, ctx context.Context, sender, receiver, classId string) {
	resp, err := ms.TransferClass(ctx, &types.MsgTransferClass{
		Creator:  sender,
		ClassId:  classId,
		Receiver: receiver,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	eventTransferClass, err := utils.GetTypedEvent(sdkCtx, &types.EventTransferClass{})

	require.NoError(t, err)
	require.NotNil(t, eventTransferClass)
	require.Equal(t, classId, eventTransferClass.Id)
	require.Equal(t, sender, eventTransferClass.OldOwner)
	require.Equal(t, receiver, eventTransferClass.NewOwner)
}

func TestTransferClass(t *testing.T) {
	ms, ctx, ctrl, nftKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	classId := mustCreateClass(t, ms, ctx, ctrl, nftKeeper, alice)

	mustTransferClass(t, ms, ctx, alice, bob, classId)
}

func TestTransferClassNotFound(t *testing.T) {
	ms, ctx, ctrl, _ := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	resp, err := ms.TransferClass(ctx, &types.MsgTransferClass{
		Creator:  alice,
		ClassId:  "1",
		Receiver: bob,
	})

	require.Nil(t, resp)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrNotFound)
}

func TestTransferClassUnauthorized(t *testing.T) {
	ms, ctx, ctrl, nftKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	classId := mustCreateClass(t, ms, ctx, ctrl, nftKeeper, alice)

	resp, err := ms.TransferClass(ctx, &types.MsgTransferClass{
		Creator:  bob,
		ClassId:  classId,
		Receiver: carol,
	})

	require.Nil(t, resp)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrUnauthorized)
}
