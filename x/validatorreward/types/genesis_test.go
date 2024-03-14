package types_test

import (
	"testing"

	"github.com/cometbft/cometbft/types/time"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/validatorreward/types"
)

func TestGenesisState_Validate(t *testing.T) {
	utils.InitSDKConfig()

	now := time.Now()

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
				Params:             types.DefaultParams(),
				LastDistributeTime: &now,
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
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
