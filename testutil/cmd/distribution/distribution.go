package distribution

import (
	"testing"

	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
)

type Pool struct {
	Pool testutil.Coins `json:"pool"`
}

func MustGetCommunityPool(t testing.TB) Pool {
	var pool Pool
	cmd.MustQuery(t, &pool, "distribution", "community-pool")
	return pool
}
