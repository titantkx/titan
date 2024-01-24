package testutil_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/utils"
)

func TestMustParseAmount(t *testing.T) {
	coins := testutil.MustParseAmount(t, "100tkx,0.5eth")

	require.Len(t, coins, 2)
	require.Equal(t, "100", coins[0].Amount.String())
	require.Equal(t, "tkx", coins[0].Denom)
	require.Equal(t, "0.5", coins[1].Amount.String())
	require.Equal(t, "eth", coins[1].Denom)
}

func TestMustGetBaseDenomAmount(t *testing.T) {
	tests := []struct {
		Amount   string
		Expected testutil.Int
	}{
		{
			Amount:   fmt.Sprintf("100%s", utils.DisplayDenom),
			Expected: testutil.MakeIntFromString("100000000000000000000"),
		},
		{
			Amount:   fmt.Sprintf("0.5%s", utils.DisplayDenom),
			Expected: testutil.MakeIntFromString("500000000000000000"),
		},
		{
			Amount:   fmt.Sprintf("300%s", utils.BaseDenom),
			Expected: testutil.MakeIntFromString("300"),
		},
		{
			Amount:   "1eth",
			Expected: testutil.MakeIntFromString("0"),
		},
		{
			Amount:   fmt.Sprintf("10%s,1eth", utils.BaseDenom),
			Expected: testutil.MakeIntFromString("10"),
		},
		{
			Amount:   fmt.Sprintf("0.1%s,1eth", utils.DisplayDenom),
			Expected: testutil.MakeIntFromString("100000000000000000"),
		},
	}

	for _, test := range tests {
		actual := testutil.MustGetBaseDenomAmount(t, test.Amount)
		require.Equal(t, test.Expected, actual)
	}
}
