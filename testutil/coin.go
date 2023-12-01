package testutil

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var coinPattern = regexp.MustCompile("((?:[\\d]+\\.)?[\\d]+)([\\w]+)")

type Coin struct {
	Amount BigFloat
	Denom  string
}

func (c *Coin) UnmarshalText(txt []byte) error {
	matches := coinPattern.FindStringSubmatch(string(txt))
	if len(matches) != 3 {
		return fmt.Errorf("invalid coin: %s", string(txt))
	}
	err := c.Amount.UnmarshalText([]byte(matches[1]))
	if err != nil {
		return fmt.Errorf("invalid coin: %s", string(txt))
	}
	c.Denom = matches[2]
	return nil
}

func MustParseAmount(t testing.TB, amount string) []Coin {
	var coins []Coin
	for _, txt := range strings.Split(amount, ",") {
		var coin Coin
		err := coin.UnmarshalText([]byte(txt))
		require.NoError(t, err)
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
