package testutil

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tokenize-titan/titan/utils"
)

var coinPattern = regexp.MustCompile(`((?:[\d]+\.)?[\d]+)([\w]+)`)

type Coin struct {
	Amount BigFloat `json:"amount"`
	Denom  string   `json:"denom"`
}

type Coins []Coin

func (coins Coins) GetBaseDenomAmount() BigInt {
	baseDenomAmount := MakeBigInt(0)
	for _, coin := range coins {
		if coin.Denom == "tkx" {
			baseDenomAmount = baseDenomAmount.Add(coin.Amount.Mul(MakeBigFloat(1000_000_000_000_000_000)).BigInt())
		} else if coin.Denom == utils.BaseDenom {
			baseDenomAmount = baseDenomAmount.Add(coin.Amount.BigInt())
		}
	}
	return baseDenomAmount
}

func ParseCoin(txt string) (Coin, error) {
	var coin Coin
	matches := coinPattern.FindStringSubmatch(txt)
	if len(matches) != 3 {
		return coin, errors.New("invalid coin format")
	}
	err := coin.Amount.UnmarshalText([]byte(matches[1]))
	if err != nil {
		return coin, err
	}
	coin.Denom = matches[2]
	return coin, nil
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

func ParseAmount(txt string) (Coins, error) {
	var coins Coins
	for _, txt := range strings.Split(txt, ",") {
		coin, err := ParseCoin(txt)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}
	return coins, nil
}

func MustParseAmount(t testing.TB, amount string) Coins {
	var coins Coins
	for _, txt := range strings.Split(amount, ",") {
		coin := MustParseCoin(t, txt)
		coins = append(coins, coin)
	}
	return coins
}

func MustGetBaseDenomAmount(t testing.TB, amount string) BigInt {
	return MustParseAmount(t, amount).GetBaseDenomAmount()
}
