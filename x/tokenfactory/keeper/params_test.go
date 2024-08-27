package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	testkeeper "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.TokenfactoryKeeper(t)
	params := types.Params{
		DenomCreationFee:        sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(1))),
		DenomCreationGasConsume: uint64(types.DefaultCreationGasFee),
	}

	err := k.SetParams(ctx, params)

	require.NoError(t, err)
	require.EqualValues(t, params, k.GetParams(ctx))
}
