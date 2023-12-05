package distribution

import (
	"testing"

	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
)

type Pool struct {
	Pool []testutil.Coin `json:"pool"`
}

func (p Pool) GetUtkxAmount() testutil.BigInt {
	for _, coin := range p.Pool {
		if coin.Denom == "utkx" {
			return coin.Amount.BigInt()
		}
	}
	return testutil.MakeBigInt(0)
}

func MustGetCommunityPool(t testing.TB) Pool {
	var pool Pool
	cmd.MustQuery(t, &pool, "distribution", "community-pool")
	return pool
}
