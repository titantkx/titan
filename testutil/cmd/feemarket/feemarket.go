package feemarket

import (
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
)

type Params struct {
	BaseFee                  testutil.Int   `json:"base_fee"`
	BaseFeeChangeDenominator int64          `json:"base_fee_change_denominator"`
	ElasticityMultiplier     int64          `json:"elasticity_multiplier"`
	EnableHeight             testutil.Int   `json:"enable_height"`
	MinGasMultiplier         testutil.Float `json:"min_gas_multiplier"`
	MinGasPrice              testutil.Float `json:"min_gas_price"`
	NoBaseFee                bool           `json:"no_base_fee"`
}

func MustGetParams(t testutil.TestingT) Params {
	var v struct {
		Params Params `json:"params"`
	}
	cmd.MustQuery(t, &v, "feemarket", "params")
	return v.Params
}

func GetBaseFee(height int64) (testutil.Int, error) {
	args := []string{
		"feemarket",
		"base-fee",
	}
	if height > 0 {
		args = append(args, "--height="+testutil.FormatInt(height))
	}
	var data struct {
		BaseFee testutil.Int `json:"base_fee"`
	}
	if err := cmd.Query(&data, args...); err != nil {
		return testutil.Int{}, err
	}
	return data.BaseFee, nil
}
