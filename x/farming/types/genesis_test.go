package types_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/x/farming/types"
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
				FarmList: []types.Farm{
					{
						Token: "tkx",
					},
					{
						Token: "btc",
					},
				},
				StakingInfoList: []types.StakingInfo{
					{
						Token:  "tkx",
						Staker: sample.AccAddress().String(),
						Amount: math.NewInt(100),
					},
					{
						Token:  "btc",
						Staker: sample.AccAddress().String(),
						Amount: math.NewInt(1000),
					},
				},
				DistributionInfo: &types.DistributionInfo{
					LastDistributionTime: time.Now(),
				},
				RewardList: []types.Reward{
					{
						Farmer: "0",
					},
					{
						Farmer: "1",
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated farm",
			genState: &types.GenesisState{
				FarmList: []types.Farm{
					{
						Token: "tkx",
					},
					{
						Token: "tkx",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated stakingInfo",
			genState: &types.GenesisState{
				StakingInfoList: []types.StakingInfo{
					{
						Token:  "tkx",
						Staker: "titan184hrxl49ktwu2aqsvjfafy47cvh850e9qanagj",
					},
					{
						Token:  "tkx",
						Staker: "titan184hrxl49ktwu2aqsvjfafy47cvh850e9qanagj",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated reward",
			genState: &types.GenesisState{
				RewardList: []types.Reward{
					{
						Farmer: "0",
					},
					{
						Farmer: "0",
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
