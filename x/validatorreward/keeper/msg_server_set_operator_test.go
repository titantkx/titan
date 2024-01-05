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

func TestMsgServer_SetOperator(t *testing.T) {
	utils.InitSDKConfig()

	zeroAddr, err := sdk.AccAddressFromHexUnsafe(zeroAddrHexStr)
	require.NoError(t, err)
	zeroAddrStr := zeroAddr.String()

	newAddrStr := sample.AccAddress()

	// Set up test cases
	testCases := []struct {
		name        string
		operator    string
		newOperator string
		expect      func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, operator string, newOperator string)
	}{
		{
			name:        "Invalid operator",
			operator:    "invalid operator address format",
			newOperator: newAddrStr,
			expect: func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, operator string, newOperator string) {
				oldOperator := k.GetOperator(ctx)
				// Call the function to be tested
				msg := types.NewMsgSetOperator(operator, newOperator)
				_, err := ms.SetOperator(ctx, msg)

				require.Error(t, err)
				require.Equal(t, oldOperator, k.GetOperator(ctx))
			},
		},
		{
			name:        "Invalid operator address",
			operator:    sample.AccAddress(),
			newOperator: newAddrStr,
			expect: func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, operator string, newOperator string) {
				oldOperator := k.GetOperator(ctx)
				// Call the function to be tested
				msg := types.NewMsgSetOperator(operator, newOperator)
				_, err := ms.SetOperator(ctx, msg)

				require.Error(t, err)
				require.Equal(t, oldOperator, k.GetOperator(ctx))
			},
		},
		{
			name:        "Valid operator",
			operator:    zeroAddrStr,
			newOperator: newAddrStr,
			expect: func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, operator string, newOperator string) {
				// Call the function to be tested
				msg := types.NewMsgSetOperator(operator, newOperator)
				_, err := ms.SetOperator(ctx, msg)

				// Check if the operator is set correctly
				require.NoError(t, err)
				require.Equal(t, newOperator, k.GetOperator(ctx).String())

				// Check if the event is emitted
				require.Equal(t, 1, len(ctx.EventManager().Events()))
				require.Equal(t, types.EventTypeSetOperator, ctx.EventManager().Events()[0].Type)
				require.Equal(t, newOperator, ctx.EventManager().Events()[0].Attributes[0].Value)
			},
		},
		{
			name:        "Invalid new operator address",
			operator:    zeroAddrStr,
			newOperator: "invalid new operator address format",
			expect: func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, operator string, newOperator string) {
				oldOperator := k.GetOperator(ctx)
				// Call the function to be tested
				msg := types.NewMsgSetOperator(operator, newOperator)
				_, err := ms.SetOperator(ctx, msg)

				require.Error(t, err)
				require.Equal(t, oldOperator, k.GetOperator(ctx))
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		ms, ctx, k := setupMsgServer(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.expect(ms, ctx, k, tc.operator, tc.newOperator)
		})
	}
}
