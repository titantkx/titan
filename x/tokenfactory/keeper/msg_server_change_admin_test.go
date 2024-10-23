package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func TestChangeAdmin(t *testing.T) {
	ms, ctx, ctrl, _, bk, _, _ := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	creator, denom := mustCreateDenom(t, ms, ctx, bk)
	newAdmin := sample.AccAddress().String()

	resp, err := ms.ChangeAdmin(ctx, &types.MsgChangeAdmin{
		Sender:   creator,
		Denom:    denom,
		NewAdmin: newAdmin,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	eventNewAdmin := utils.GetABCIEvent(sdkCtx, types.TypeMsgChangeAdmin)

	require.NotNil(t, eventNewAdmin)
	require.Equal(t, denom, utils.GetABCIEventAttribute(eventNewAdmin, types.AttributeDenom))
	require.Equal(t, newAdmin, utils.GetABCIEventAttribute(eventNewAdmin, types.AttributeNewAdmin))
}
