package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func TestMint(t *testing.T) {
	ms, ctx, ctrl, _, bk, _, _ := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	creator, denom := mustCreateDenom(t, ms, ctx, bk)
	amount := sdk.NewCoin(denom, sdk.NewInt(1000))
	receiver := sample.AccAddress()

	bk.EXPECT().GetDenomMetaData(ctx, denom).Return(banktypes.Metadata{}, true)
	bk.EXPECT().MintCoins(ctx, types.ModuleName, sdk.NewCoins(amount)).Return(nil)
	bk.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(amount)).Return(nil)

	resp, err := ms.Mint(ctx, &types.MsgMint{
		Sender:        creator,
		Amount:        amount,
		MintToAddress: receiver.String(),
	})

	require.NoError(t, err)
	require.NotNil(t, resp)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	eventMint := utils.GetABCIEvent(sdkCtx, types.TypeMsgMint)

	require.NotNil(t, eventMint)
	require.Equal(t, receiver.String(), utils.GetABCIEventAttribute(eventMint, types.AttributeMintToAddress))
	require.Equal(t, amount.String(), utils.GetABCIEventAttribute(eventMint, types.AttributeAmount))
}
