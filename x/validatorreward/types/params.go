package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyRate = []byte("Rate")
	// TODO: Determine the default value
	DefaultRate string = "rate"
)

var (
	KeyOperator = []byte("Operator")
	// TODO: Determine the default value
	DefaultOperator string = "operator"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	rate string,
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

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyRate, &p.Rate, validateRate),
		paramtypes.NewParamSetPair(KeyOperator, &p.Operator, validateOperator),
	}
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
	rate, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = rate

	return nil
}

// validateOperator validates the Operator param
func validateOperator(v interface{}) error {
	operator, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = operator

	return nil
}
