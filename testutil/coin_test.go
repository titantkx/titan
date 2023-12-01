package testutil_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
)

func TestMustParseAmount(t *testing.T) {
	coins := testutil.MustParseAmount(t, "100tkx,0.5eth")

	require.Len(t, coins, 2)
	require.Equal(t, testutil.MakeBigFloat(100), coins[0].Amount)
	require.Equal(t, "tkx", coins[0].Denom)
	require.Equal(t, testutil.MakeBigFloat(0.5), coins[1].Amount)
	require.Equal(t, "eth", coins[1].Denom)
}

func TestMustGetUtkxAmount(t *testing.T) {
	tests := []struct {
		Amount   string
		Expected testutil.BigInt
	}{
		{
			Amount:   "100tkx",
			Expected: testutil.MakeBigIntFromString("100000000000000000000"),
		},
		{
			Amount:   "0.5tkx",
			Expected: testutil.MakeBigIntFromString("500000000000000000"),
		},
		{
			Amount:   "300utkx",
			Expected: testutil.MakeBigIntFromString("300"),
		},
		{
			Amount:   "1eth",
			Expected: testutil.MakeBigIntFromString("0"),
		},
		{
			Amount:   "10utkx,1eth",
			Expected: testutil.MakeBigIntFromString("10"),
		},
		{
			Amount:   "0.1tkx,1eth",
			Expected: testutil.MakeBigIntFromString("100000000000000000"),
		},
	}

	for _, test := range tests {
		actual := testutil.MustGetUtkxAmount(t, test.Amount)
		require.Equal(t, test.Expected, actual)
	}
}
