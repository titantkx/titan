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

func (coin Coin) String() string {
	return coin.Amount.String() + coin.Denom
}

func (coin Coin) GetBaseDenomAmount() BigInt {
	switch coin.Denom {
	case utils.DisplayDenom:
		return coin.Amount.Mul(MakeBigFloat(1_000_000_000_000_000_000)).BigInt()
	case utils.MilliDenom:
		return coin.Amount.Mul(MakeBigFloat(1_000_000_000_000_000)).BigInt()
	case utils.MicroDenom:
		return coin.Amount.Mul(MakeBigFloat(1_000_000_000_000)).BigInt()
	case utils.BaseDenom:
		return coin.Amount.BigInt()
	default:
		return MakeBigInt(0)
	}
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

type Coins []Coin

func (coins Coins) String() string {
	s := make([]string, len(coins))
	for i := range coins {
		s[i] = coins[i].String()
	}
	return strings.Join(s, ",")
}

func (coins Coins) GetBaseDenomAmount() BigInt {
	baseDenomAmount := MakeBigInt(0)
	for _, coin := range coins {
		baseDenomAmount = baseDenomAmount.Add(coin.GetBaseDenomAmount())
	}
	return baseDenomAmount
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
