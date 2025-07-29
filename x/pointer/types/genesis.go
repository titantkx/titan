package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Erc20NativeList: []Erc20Native{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in erc20Native
	erc20NativeIndexMap := make(map[string]struct{})

	for _, elem := range gs.Erc20NativeList {
		index := string(Erc20NativeKey(elem.TokenDenom))
		if _, ok := erc20NativeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for erc20Native")
		}
		erc20NativeIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
