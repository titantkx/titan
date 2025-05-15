package pointer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/pointer"
	"github.com/titantkx/titan/x/pointer/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		Erc20NativeList: []types.Erc20Native{
			{
				TokenDenom: "0",
			},
			{
				TokenDenom: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.PointerKeeper(t)
	pointer.InitGenesis(ctx, *k, genesisState)
	got := pointer.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.Erc20NativeList, got.Erc20NativeList)
	// this line is used by starport scaffolding # genesis/test/assert
}
