package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/farming/types"
)

func TestStake(t *testing.T) {
	ms, ctx, ctrl, k, bankKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	msg := &types.MsgStake{
		Sender: sample.AccAddress().String(),
		Amount: utils.NewCoins("1000tkx"),
	}

	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(msg.Sender), types.ModuleName, msg.Amount).Return(nil)

	resp, err := ms.Stake(ctx, msg)

	require.NoError(t, err)
	require.NotNil(t, resp)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	eventStake, err := utils.GetTypedEvent(sdkCtx, &types.EventStake{})

	require.NoError(t, err)
	require.NotNil(t, eventStake)
	require.Equal(t, msg.Sender, eventStake.Sender)
	require.Equal(t, msg.Amount, eventStake.Amount)

	stakingInfo, found := k.GetStakingInfo(sdkCtx, "tkx", msg.Sender)

	require.True(t, found)
	require.Equal(t, sdk.NewInt(1000), stakingInfo.Amount)
}
