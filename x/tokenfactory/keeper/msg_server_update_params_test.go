package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func TestUpdateParams(t *testing.T) {
	ms, ctx := setupMsgServer(t)

	testCases := []struct {
		name      string
		input     types.MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "valid authority",
			input: types.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: types.Params{
					DenomCreationFee:        []sdk.Coin{},
					DenomCreationGasConsume: 0,
				},
			},
			expErr: false,
		},
		{
			name: "invalid authority",
			input: types.MsgUpdateParams{
				Authority: sample.AccAddress().String(),
				Params: types.Params{
					DenomCreationFee:        []sdk.Coin{},
					DenomCreationGasConsume: 0,
				},
			},
			expErr:    true,
			expErrMsg: "unauthorized account",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ms.UpdateParams(ctx, &tc.input)

			if tc.expErr {
				require.Nil(t, resp)
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expErrMsg)
			} else {
				require.NotNil(t, resp)
				require.NoError(t, err)

			}
		})
	}
}
