package distribution

import (
	"testing"

	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
)

type Params struct {
	CommunityTax        testutil.Float `json:"community_tax"`
	WithdrawAddrEnabled bool           `json:"withdraw_addr_enabled"`
}

func MustGetParams(t testing.TB) Params {
	var v Params
	cmd.MustQuery(t, &v, "distribution", "params")
	return v
}

type Pool struct {
	Pool testutil.Coins `json:"pool"`
}

func MustGetCommunityPool(t testing.TB) Pool {
	var pool Pool
	cmd.MustQuery(t, &pool, "distribution", "community-pool")
	return pool
}

func MustGetRewards(t testing.TB, del string, val string, height int64) testutil.Coins {
	args := []string{"distribution", "rewards", del}
	if val != "" {
		args = append(args, val)
	}
	if height > 0 {
		args = append(args, "--height="+testutil.FormatInt(height))
	}
	if val == "" {
		var v struct {
			Total testutil.Coins `json:"total"`
		}
		cmd.MustQuery(t, &v, args...)
		return v.Total
	} else {
		var v struct {
			Rewards testutil.Coins `json:"rewards"`
		}
		cmd.MustQuery(t, &v, args...)
		return v.Rewards
	}
}
