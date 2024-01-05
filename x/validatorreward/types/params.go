package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var DefaultRate = sdk.NewDecWithPrec(6, 2) // 6%

var DefaultAuthority string = "titan1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqe22t5s"

// NewParams creates a new Params instance
func NewParams(
	rate sdkmath.LegacyDec,
	authority string,
) Params {
	return Params{
		Rate:      rate,
		Authority: authority,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultRate,
		DefaultAuthority,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateRate(p.Rate); err != nil {
		return err
	}

	if err := validateAuthority(p.Authority); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// validateRate validates the Rate param
func validateRate(v interface{}) error {
	rate, ok := v.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if rate.IsNil() {
		return fmt.Errorf("rate should not be nil")
	}

	if rate.IsNegative() {
		return fmt.Errorf("rate should not be negative")
	}

	if rate.GT(sdk.OneDec()) {
		return fmt.Errorf("rate should not be greater than 1")
	}

	return nil
}

// validateAuthority validates the Authority param
func validateAuthority(v interface{}) error {
	authority, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// validate authority address
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		return fmt.Errorf("invalid authority address: %s", err)
	}

	return nil
}
