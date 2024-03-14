package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/validatorreward/keeper"
	"github.com/titantkx/titan/x/validatorreward/types"
)

func TestMsgServer_SetAuthority(t *testing.T) {
	utils.InitSDKConfig()

	zeroAddr, err := sdk.AccAddressFromHexUnsafe(types.ZeroHexAddress)
	require.NoError(t, err)
	zeroAddrStr := zeroAddr.String()

	newAddrStr := sample.AccAddress()

	// Set up test cases
	testCases := []struct {
		name         string
		authority    string
		newAuthority string
		expect       func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, authority string, newAuthority string)
	}{
		{
			name:         "Invalid authority",
			authority:    "invalid authority address format",
			newAuthority: newAddrStr,
			expect: func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, authority string, newAuthority string) {
				oldAuthority := k.GetAuthority(ctx)
				// Call the function to be tested
				msg := types.NewMsgSetAuthority(authority, newAuthority)
				_, err := ms.SetAuthority(ctx, msg)

				require.Error(t, err)
				require.Equal(t, oldAuthority, k.GetAuthority(ctx))
			},
		},
		{
			name:         "Invalid authority address",
			authority:    sample.AccAddress(),
			newAuthority: newAddrStr,
			expect: func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, authority string, newAuthority string) {
				oldAuthority := k.GetAuthority(ctx)
				// Call the function to be tested
				msg := types.NewMsgSetAuthority(authority, newAuthority)
				_, err := ms.SetAuthority(ctx, msg)

				require.Error(t, err)
				require.Equal(t, oldAuthority, k.GetAuthority(ctx))
			},
		},
		{
			name:         "Valid authority",
			authority:    zeroAddrStr,
			newAuthority: newAddrStr,
			expect: func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, authority string, newAuthority string) {
				// Call the function to be tested
				msg := types.NewMsgSetAuthority(authority, newAuthority)
				_, err := ms.SetAuthority(ctx, msg)

				// Check if the authority is set correctly
				require.NoError(t, err)
				require.Equal(t, newAuthority, k.GetAuthority(ctx).String())

				// Check if the event is emitted
				require.Equal(t, 1, len(ctx.EventManager().Events()))
				require.Equal(t, types.EventTypeSetAuthority, ctx.EventManager().Events()[0].Type)
				require.Equal(t, newAuthority, ctx.EventManager().Events()[0].Attributes[0].Value)
			},
		},
		{
			name:         "Invalid new authority address",
			authority:    zeroAddrStr,
			newAuthority: "invalid new authority address format",
			expect: func(ms types.MsgServer, ctx sdk.Context, k *keeper.Keeper, authority string, newAuthority string) {
				oldAuthority := k.GetAuthority(ctx)
				// Call the function to be tested
				msg := types.NewMsgSetAuthority(authority, newAuthority)
				_, err := ms.SetAuthority(ctx, msg)

				require.Error(t, err)
				require.Equal(t, oldAuthority, k.GetAuthority(ctx))
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		ms, ctx, k := setupMsgServer(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.expect(ms, ctx, k, tc.authority, tc.newAuthority)
		})
	}
}
