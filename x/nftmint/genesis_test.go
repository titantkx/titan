package nftmint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/nftmint"
	"github.com/titantkx/titan/x/nftmint/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		SystemInfo: types.SystemInfo{
			NextClassId: 8,
		},
		MintingInfoList: []types.MintingInfo{
			{
				ClassId: "0",
			},
			{
				ClassId: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.NftmintKeeper(t)
	nftmint.InitGenesis(ctx, *k, genesisState)
	got := nftmint.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.MintingInfoList, got.MintingInfoList)
	require.Equal(t, genesisState.SystemInfo, got.SystemInfo)
	// this line is used by starport scaffolding # genesis/test/assert
}
