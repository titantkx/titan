package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var DefaultRate = sdk.NewDecWithPrec(6, 2) // 6%

var DefaultOperator string = "titan1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqe22t5s"

// NewParams creates a new Params instance
func NewParams(
	rate sdkmath.LegacyDec,
	operator string,
) Params {
	return Params{
		Rate:     rate,
		Operator: operator,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultRate,
		DefaultOperator,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateRate(p.Rate); err != nil {
		return err
	}

	if err := validateOperator(p.Operator); err != nil {
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

// validateOperator validates the Operator param
func validateOperator(v interface{}) error {
	operator, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// validate operator address
	if _, err := sdk.AccAddressFromBech32(operator); err != nil {
		return fmt.Errorf("invalid operator address: %s", err)
	}

	return nil
}
