package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tokenize-titan/titan/testutil/sample"
	"github.com/tokenize-titan/titan/utils"
	"github.com/tokenize-titan/titan/x/validatorreward/keeper"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

func TestMsgServer_SetRate(t *testing.T) {
	utils.InitSDKConfig()

	zeroAddr, err := sdk.AccAddressFromHexUnsafe(zeroAddrHexStr)
	require.NoError(t, err)
	zeroAddrStr := zeroAddr.String()

	// Set up test cases
	testCases := []struct {
		name     string
		operator string
		rate     sdk.Dec
		expect   func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, operator string, rate sdk.Dec)
	}{
		{
			name:     "Valid operator",
			operator: zeroAddrStr,
			rate:     sdk.NewDecWithPrec(5, 1),
			expect: func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, operator string, rate sdk.Dec) {
				// Call the function to be tested
				msg := types.NewMsgSetRate(operator, rate)
				_, err := ms.SetRate(ctx, msg)

				// Check if the rate is set correctly
				require.NoError(t, err)
				require.Equal(t, rate, k.GetRate(ctx))
			},
		},
		{
			name:     "Invalid operator",
			operator: "invalid operator address format",
			rate:     sdk.NewDecWithPrec(5, 1),
			expect: func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, operator string, rate sdk.Dec) {
				oldRate := k.GetRate(ctx)
				// Call the function to be tested
				msg := types.NewMsgSetRate(operator, rate)
				_, err := ms.SetRate(ctx, msg)

				require.Error(t, err)
				require.Equal(t, oldRate, k.GetRate(ctx))
			},
		},
		{
			name:     "Invalid operator address",
			operator: sample.AccAddress(),
			rate:     sdk.NewDecWithPrec(5, 1),
			expect: func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, operator string, rate sdk.Dec) {
				oldRate := k.GetRate(ctx)
				// Call the function to be tested
				msg := types.NewMsgSetRate(operator, rate)
				_, err := ms.SetRate(ctx, msg)

				require.ErrorIs(t, err, types.ErrForbidden)
				require.Equal(t, oldRate, k.GetRate(ctx))
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		ms, ctx, k := setupMsgServer(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.expect(ms, ctx, k, tc.operator, tc.rate)
		})
	}
}
