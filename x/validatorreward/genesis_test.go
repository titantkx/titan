package validatorreward_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/tokenize-titan/titan/testutil/keeper"
	"github.com/tokenize-titan/titan/testutil/nullify"
	"github.com/tokenize-titan/titan/x/validatorreward"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ValidatorrewardKeeper(t)
	validatorreward.InitGenesis(ctx, *k, genesisState)
	got := validatorreward.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
