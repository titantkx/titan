package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/farming/types"
)

func TestHarvest(t *testing.T) {
	ms, ctx, ctrl, k, bankKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sender := sample.AccAddress()
	reward := utils.NewCoins("1000tkx")

	k.SetReward(sdkCtx, types.Reward{
		Farmer: sender.String(),
		Amount: reward,
	})

	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, reward).Return(nil)

	resp, err := ms.Harvest(ctx, &types.MsgHarvest{
		Sender: sender.String(),
	})

	require.NoError(t, err)
	require.NotNil(t, resp)

	eventHarvest, err := utils.GetTypedEvent(sdkCtx, &types.EventHarvest{})

	require.NoError(t, err)
	require.NotNil(t, eventHarvest)
	require.Equal(t, sender.String(), eventHarvest.Sender)
	require.Equal(t, reward, eventHarvest.Amount)

	_, found := k.GetReward(sdkCtx, sender.String())

	require.False(t, found)
}
