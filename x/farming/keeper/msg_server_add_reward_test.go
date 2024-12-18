package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/farming/types"
)

func TestAddReward(t *testing.T) {
	ms, ctx, ctrl, k, bankKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	msg := &types.MsgAddReward{
		Sender:    sample.AccAddress().String(),
		Token:     "bitcoin",
		Amount:    sdk.NewCoins(sdk.NewCoin("tkx", sdk.NewInt(1000))),
		EndTime:   time.Now().Add(300 * time.Hour),
		StartTime: time.Now().Add(1 * time.Hour),
	}

	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(msg.Sender), types.ModuleName, msg.Amount).Return(nil)

	resp, err := ms.AddReward(ctx, msg)

	require.NoError(t, err)
	require.NotNil(t, resp)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	eventAddReward, err := utils.GetTypedEvent(sdkCtx, &types.EventAddReward{})

	require.NoError(t, err)
	require.NotNil(t, eventAddReward)
	require.Equal(t, msg.Sender, eventAddReward.Sender)
	require.Equal(t, msg.Token, eventAddReward.Token)
	require.Equal(t, msg.Amount, eventAddReward.Amount)
	require.True(t, msg.StartTime.Equal(eventAddReward.StartTime))
	require.True(t, msg.EndTime.Equal(eventAddReward.EndTime))

	farm, found := k.GetFarm(sdkCtx, "bitcoin")

	require.True(t, found)
	require.NotEmpty(t, farm.Rewards)

	reward := farm.Rewards[0]

	require.Equal(t, msg.Sender, reward.Sender)
	require.Equal(t, msg.Amount, reward.Amount)
	require.True(t, msg.StartTime.Equal(eventAddReward.StartTime))
	require.True(t, msg.EndTime.Equal(eventAddReward.EndTime))
}
