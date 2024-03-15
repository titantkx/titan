package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/x/nftmint/testutil"
	"github.com/titantkx/titan/x/nftmint/types"
)

func msgMint(minter string, receiver string, classId string) *types.MsgMint {
	return &types.MsgMint{
		Creator:  minter,
		Receiver: receiver,
		ClassId:  classId,
		Uri:      sample.URL(),
		UriHash:  sample.Hash(),
		Data:     sample.JSON(),
	}
}

func mustMintNFT(t testing.TB, ms types.MsgServer, ctx context.Context, ctrl *gomock.Controller, nftKeeper *testutil.MockNFTKeeper, minter string, receiver string, classId string) string {
	nftKeeper.EXPECT().Mint(ctx, gomock.Any(), sdk.MustAccAddressFromBech32(receiver)).Return(nil)

	resp, err := ms.Mint(ctx, msgMint(minter, receiver, classId))

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Id)

	return resp.Id
}

func TestMint(t *testing.T) {
	ms, ctx, ctrl, nftKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	classId := mustCreateClass(t, ms, ctx, ctrl, nftKeeper, alice)

	mustMintNFT(t, ms, ctx, ctrl, nftKeeper, alice, bob, classId)
}

func TestMintClassNotFound(t *testing.T) {
	ms, ctx, ctrl, _ := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	resp, err := ms.Mint(ctx, msgMint(alice, bob, "1"))

	require.Nil(t, resp)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrNotFound)
}

func TestMintUnauthorized(t *testing.T) {
	ms, ctx, ctrl, nftKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	classId := mustCreateClass(t, ms, ctx, ctrl, nftKeeper, alice)

	resp, err := ms.Mint(ctx, msgMint(carol, bob, classId))

	require.Nil(t, resp)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrUnauthorized)
}
