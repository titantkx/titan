package keeper_test

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/tokenfactory/testutil"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

//nolint:revive
func mustCreateDenom(t *testing.T, ms types.MsgServer, ctx context.Context, bk *testutil.MockBankKeeper) (string, string) {
	creator := sample.AccAddress().String()
	subDenom := "bitcoin"
	denom := fmt.Sprintf("factory/%s/%s", creator, subDenom)

	bk.EXPECT().HasSupply(ctx, subDenom).Return(false)
	bk.EXPECT().GetDenomMetaData(ctx, denom).Return(banktypes.Metadata{}, false)
	bk.EXPECT().GetDenomMetaData(ctx, denom).Return(banktypes.Metadata{}, false)
	bk.EXPECT().SetDenomMetaData(ctx, banktypes.Metadata{
		DenomUnits: []*banktypes.DenomUnit{{
			Denom:    denom,
			Exponent: 0,
		}},
		Base:    denom,
		Name:    denom,
		Symbol:  denom,
		Display: denom,
	})

	resp, err := ms.CreateDenom(ctx, &types.MsgCreateDenom{
		Sender:   creator,
		Subdenom: subDenom,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, denom, resp.NewTokenDenom)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	eventCreateDenom := utils.GetABCIEvent(sdkCtx, types.TypeMsgCreateDenom)

	require.NotNil(t, eventCreateDenom)
	require.Equal(t, creator, utils.GetABCIEventAttribute(eventCreateDenom, types.AttributeCreator))
	require.Equal(t, denom, utils.GetABCIEventAttribute(eventCreateDenom, types.AttributeNewTokenDenom))

	return creator, denom
}

func TestCreateDenom(t *testing.T) {
	ms, ctx, ctrl, _, bk, _, _ := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	mustCreateDenom(t, ms, ctx, bk)
}
