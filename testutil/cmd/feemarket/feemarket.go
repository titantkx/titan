package feemarket

import (
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
)

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
