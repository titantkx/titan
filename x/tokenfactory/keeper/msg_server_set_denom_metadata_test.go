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

func TestSetDenomMetadata(t *testing.T) {
	ms, ctx, ctrl, _, bk, _, _ := setupMsgServerWithMocks(t)
	defer ctrl.Finish()

	creator, denom := mustCreateDenom(t, ms, ctx, bk)
	metadata := banktypes.Metadata{
		Description: sample.Sentence(),
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    denom,
				Exponent: 0,
			},
			{
				Denom:    "tkx",
				Exponent: 18,
			},
		},
		Base:    denom,
		Display: "tkx",
		Name:    sample.Name(),
		Symbol:  sample.Word(),
		URI:     sample.URL(),
		URIHash: sample.Hash(),
	}

	bk.EXPECT().SetDenomMetaData(ctx, metadata)

	resp, err := ms.SetDenomMetadata(ctx, &types.MsgSetDenomMetadata{
		Sender:   creator,
		Metadata: metadata,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	eventSetDenomMetadata := utils.GetABCIEvent(sdkCtx, types.TypeMsgSetDenomMetadata)

	require.NotNil(t, eventSetDenomMetadata)
	require.Equal(t, denom, utils.GetABCIEventAttribute(eventSetDenomMetadata, types.AttributeDenom))
	require.Equal(t, metadata.String(), utils.GetABCIEventAttribute(eventSetDenomMetadata, types.AttributeDenomMetadata))
}
