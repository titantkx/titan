package feemarket

import (
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
)

func GetBaseFee(height int64) (testutil.BigInt, error) {
	args := []string{
		"feemarket",
		"base-fee",
	}
	if height > 0 {
		args = append(args, "--height="+testutil.FormatInt(height))
	}
	var data struct {
		BaseFee testutil.BigInt `json:"base_fee"`
	}
	if err := cmd.Query(&data, args...); err != nil {
		return testutil.BigInt{}, err
	}
	return data.BaseFee, nil
}
