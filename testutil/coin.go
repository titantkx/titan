package testutil

import (
	"errors"
	"regexp"
	"strings"

	"github.com/stretchr/testify/require"

	"github.com/titantkx/titan/utils"
)

var coinPattern = regexp.MustCompile(`((?:[\d]+\.)?[\d]+)([\w]+)`)

var OneToken Coin

func init() {
	var err error
	OneToken, err = ParseCoin("1" + utils.DisplayDenom)
	if err != nil {
		panic(err)
	}
}

type Coin struct {
	Amount Float  `json:"amount"`
	Denom  string `json:"denom"`
}

func (coin Coin) Mul(x float64) Coin {
	return Coin{
		Amount: coin.Amount.Mul(MakeFloat(x)),
		Denom:  coin.Denom,
	}
}

func (coin Coin) GetAmount() string {
	return coin.GetBaseDenomAmount().String() + utils.BaseDenom
}

func (coin Coin) String() string {
	return coin.Amount.String() + coin.Denom
}

func (coin Coin) GetBaseDenomAmount() Int {
	switch coin.Denom {
	case utils.DisplayDenom:
		return coin.Amount.Mul(MakeFloat(1_000_000_000_000_000_000)).Int()
	case utils.MilliDenom:
		return coin.Amount.Mul(MakeFloat(1_000_000_000_000_000)).Int()
	case utils.MicroDenom:
		return coin.Amount.Mul(MakeFloat(1_000_000_000_000)).Int()
	case utils.BaseDenom:
		return coin.Amount.Int()
	default:
		return MakeInt(0)
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

func MustParseCoin(t TestingT, txt string) Coin {
	var coin Coin
	matches := coinPattern.FindStringSubmatch(txt)
	require.Len(t, matches, 3)
	err := coin.Amount.UnmarshalText([]byte(matches[1]))
	require.NoError(t, err)
	coin.Denom = matches[2]
	return coin
}

type Coins []Coin

func (coins Coins) Add(coin Coin) Coins {
	return append(coins, coin)
}

func (coins Coins) GetAmount() string {
	return coins.GetBaseDenomAmount().String() + utils.BaseDenom
}

func (coins Coins) String() string {
	s := make([]string, len(coins))
	for i := range coins {
		s[i] = coins[i].String()
	}
	return strings.Join(s, ",")
}

func (coins Coins) GetBaseDenomAmount() Int {
	baseDenomAmount := MakeInt(0)
	for _, coin := range coins {
		baseDenomAmount = baseDenomAmount.Add(coin.GetBaseDenomAmount())
	}
	return baseDenomAmount
}

func ParseAmount(txt string) (Coins, error) {
	coinTxts := strings.Split(txt, ",")
	coins := make(Coins, 0, len(coinTxts))
	for _, txt := range coinTxts {
		coin, err := ParseCoin(txt)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}
	return coins, nil
}

func MustParseAmount(t TestingT, amount string) Coins {
	coinTxts := strings.Split(amount, ",")
	coins := make(Coins, 0, len(coinTxts))
	for _, txt := range coinTxts {
		coin := MustParseCoin(t, txt)
		coins = append(coins, coin)
	}
	return coins
}

func MustGetBaseDenomAmount(t TestingT, amount string) Int {
	return MustParseAmount(t, amount).GetBaseDenomAmount()
}
