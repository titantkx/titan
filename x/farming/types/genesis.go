package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		FarmList:         []Farm{},
		StakingInfoList:  []StakingInfo{},
		DistributionInfo: nil,
		RewardList:       []Reward{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in farm
	farmIndexMap := make(map[string]struct{})

	for _, elem := range gs.FarmList {
		index := string(FarmKey(elem.Token))
		if _, ok := farmIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for farm")
		}
		farmIndexMap[index] = struct{}{}

		for _, reward := range elem.Rewards {
			_, err := sdk.AccAddressFromBech32(reward.Sender)
			if err != nil {
				return WrapErrorf(sdkerrors.ErrInvalidAddress, "invalid reward sender address (%s)", err)
			}

			if !reward.Amount.IsValid() || !reward.Amount.IsAllPositive() {
				return WrapError(sdkerrors.ErrInvalidCoins, reward.Amount.String())
			}

			if reward.EndTime.IsZero() {
				return WrapErrorf(ErrInvalidTime, "reward end time cannot be zero")
			}

			if !reward.StartTime.IsZero() && !reward.StartTime.Before(reward.EndTime) {
				return WrapErrorf(ErrInvalidTime, "reward start time must be smaller than end time")
			}
		}
	}
	// Check for duplicated index in stakingInfo
	stakingInfoIndexMap := make(map[string]struct{})

	for _, elem := range gs.StakingInfoList {
		index := string(StakingInfoKey(elem.Token, elem.Staker))
		if _, ok := stakingInfoIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for stakingInfo")
		}
		stakingInfoIndexMap[index] = struct{}{}

		if sdk.ValidateDenom(elem.Token) != nil {
			return WrapErrorf(ErrInvalidToken, "invalid token: %s", elem.Token)
		}

		_, err := sdk.AccAddressFromBech32(elem.Staker)
		if err != nil {
			return WrapErrorf(sdkerrors.ErrInvalidAddress, "invalid staker address (%s)", err)
		}

		if !elem.Amount.IsPositive() {
			return ErrInvalidStakingAmount
		}
	}
	// Check for duplicated index in reward
	rewardIndexMap := make(map[string]struct{})

	for _, elem := range gs.RewardList {
		index := string(RewardKey(elem.Farmer))
		if _, ok := rewardIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for reward")
		}
		rewardIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
