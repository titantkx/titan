package testutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokenize-titan/titan/testutil"
)

func TestFloat_String(t *testing.T) {
	tests := []struct {
		input         testutil.Float
		exepectOutput string
	}{
		{
			testutil.MakeFloat(0.1),
			"0.1",
		},
		{
			testutil.MakeFloat(111),
			"111",
		},
		{
			testutil.MakeFloatFromString("111.000"),
			"111",
		},
		{
			testutil.MakeFloat(111.111),
			"111.111",
		},
		{
			testutil.MakeFloatFromString("111.111000"),
			"111.111",
		},
		{
			testutil.MakeFloatFromString("111.111111111111111111"),
			"111.111111111111111111",
		},
		{
			testutil.MakeFloatFromString("111.111111111111111111111"),
			"111.111111111111111111",
		},
		{
			testutil.MakeFloatFromString("2.5e+20"),
			"250000000000000000000",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.exepectOutput, test.input.String())
	}
}
