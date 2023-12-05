package testutil

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var coinPattern = regexp.MustCompile(`((?:[\d]+\.)?[\d]+)([\w]+)`)

type Coin struct {
	Amount BigFloat `json:"amount"`
	Denom  string   `json:"denom"`
}

func MustParseCoin(t testing.TB, txt string) Coin {
	var coin Coin
	matches := coinPattern.FindStringSubmatch(txt)
	require.Len(t, matches, 3)
	err := coin.Amount.UnmarshalText([]byte(matches[1]))
	require.NoError(t, err)
	coin.Denom = matches[2]
	return coin
}

func MustParseAmount(t testing.TB, amount string) []Coin {
	var coins []Coin
	for _, txt := range strings.Split(amount, ",") {
		coin := MustParseCoin(t, txt)
		coins = append(coins, coin)
	}
	return coins
}

func MustGetUtkxAmount(t testing.TB, amount string) BigInt {
	coins := MustParseAmount(t, amount)
	utkxAmount := MakeBigInt(0)
	for _, coin := range coins {
		if coin.Denom == "tkx" {
			utkxAmount = utkxAmount.Add(coin.Amount.Mul(MakeBigFloat(1000_000_000_000_000_000)).BigInt())
		} else if coin.Denom == "utkx" {
			utkxAmount = utkxAmount.Add(coin.Amount.BigInt())
		}
	}
	return utkxAmount
}
