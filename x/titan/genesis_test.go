package titan_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/tokenize-titan/titan/testutil/keeper"
	"github.com/tokenize-titan/titan/testutil/nullify"
	"github.com/tokenize-titan/titan/x/titan"
	"github.com/tokenize-titan/titan/x/titan/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TitanKeeper(t)
	titan.InitGenesis(ctx, *k, genesisState)
	got := titan.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
