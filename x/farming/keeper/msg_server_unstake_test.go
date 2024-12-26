package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/farming/types"
)

func TestUnstake(t *testing.T) {
	ms, ctx, ctrl, k, bankKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sender := sample.AccAddress()

	k.SetStakingInfo(sdkCtx, types.StakingInfo{
		Token:  "tkx",
		Staker: sender.String(),
		Amount: sdk.NewInt(10000),
	})

	msg := &types.MsgUnstake{
		Sender: sender.String(),
		Amount: utils.NewCoins("1000tkx"),
	}

	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, msg.Amount).Return(nil)

	resp, err := ms.Unstake(ctx, msg)

	require.NoError(t, err)
	require.NotNil(t, resp)

	eventUnstake, err := utils.GetTypedEvent(sdkCtx, &types.EventUnstake{})

	require.NoError(t, err)
	require.NotNil(t, eventUnstake)
	require.Equal(t, msg.Sender, eventUnstake.Sender)
	require.Equal(t, msg.Amount, eventUnstake.Amount)

	stakingInfo, found := k.GetStakingInfo(sdkCtx, "tkx", sender.String())

	require.True(t, found)
	require.Equal(t, sdk.NewInt(9000), stakingInfo.Amount)
}

func TestUnstakeAl(t *testing.T) {
	ms, ctx, ctrl, k, bankKeeper := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sender := sample.AccAddress()

	k.SetStakingInfo(sdkCtx, types.StakingInfo{
		Token:  "tkx",
		Staker: sender.String(),
		Amount: sdk.NewInt(1000),
	})

	msg := &types.MsgUnstake{
		Sender: sender.String(),
		Amount: utils.NewCoins("1000tkx"),
	}

	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, msg.Amount).Return(nil)

	resp, err := ms.Unstake(ctx, msg)

	require.NoError(t, err)
	require.NotNil(t, resp)

	eventUnstake, err := utils.GetTypedEvent(sdkCtx, &types.EventUnstake{})

	require.NoError(t, err)
	require.NotNil(t, eventUnstake)
	require.Equal(t, msg.Sender, eventUnstake.Sender)
	require.Equal(t, msg.Amount, eventUnstake.Amount)

	_, found := k.GetStakingInfo(sdkCtx, "tkx", sender.String())

	require.False(t, found)
}
