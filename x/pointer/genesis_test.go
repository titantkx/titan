package pointer_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	etherminttests "github.com/titantkx/ethermint/tests"

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
				Erc20Addr:  etherminttests.GenerateAddress().String(),
			},
			{
				TokenDenom: "1",
				Erc20Addr:  etherminttests.GenerateAddress().String(),
			},
			{
				TokenDenom: "2",
				// empty Erc20Addr to test ignore this record
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.PointerKeeper(t)
	pointer.InitGenesis(ctx, *k, genesisState)
	got := pointer.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// remove record that do not have Erc20Addr
	shouldImportedErc20NativeList := append([]types.Erc20Native{}, genesisState.Erc20NativeList[:2]...)
	require.ElementsMatch(t, shouldImportedErc20NativeList, got.Erc20NativeList)
	// this line is used by starport scaffolding # genesis/test/assert

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
