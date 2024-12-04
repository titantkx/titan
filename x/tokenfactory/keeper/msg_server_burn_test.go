package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func TestBurn(t *testing.T) {
	ms, ctx, ctrl, _, bk, _, _ := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	creator, denom := mustCreateDenom(t, ms, ctx, bk)
	amount := sdk.NewCoin(denom, sdk.NewInt(1000))

	bk.EXPECT().SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(creator), types.ModuleName, sdk.NewCoins(amount)).Return(nil)
	bk.EXPECT().BurnCoins(ctx, types.ModuleName, sdk.NewCoins(amount)).Return(nil)

	resp, err := ms.Burn(ctx, &types.MsgBurn{
		Sender: creator,
		Amount: amount,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	eventMint := utils.GetABCIEvent(sdkCtx, types.TypeMsgBurn)

	require.NotNil(t, eventMint)
	require.Equal(t, creator, utils.GetABCIEventAttribute(eventMint, types.AttributeBurnFromAddress))
	require.Equal(t, amount.String(), utils.GetABCIEventAttribute(eventMint, types.AttributeAmount))
}
