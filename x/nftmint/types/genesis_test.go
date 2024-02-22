package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/x/nftmint/types"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
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
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated mintingInfo",
			genState: &types.GenesisState{

				SystemInfo: types.SystemInfo{
					NextClassId: 8,
				},
				MintingInfoList: []types.MintingInfo{
					{
						ClassId: "0",
					},
					{
						ClassId: "0",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
