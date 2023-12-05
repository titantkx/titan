package slashing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
)

type Params struct {
	DowntimeJailDuration    testutil.Duration `json:"downtime_jail_duration"`
	MinSignedPerWindow      testutil.BigFloat `json:"min_signed_per_window"`
	SignedBlocksWindow      testutil.Int      `json:"signed_blocks_window"`
	SlashFractionDoubleSign testutil.BigFloat `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   testutil.BigFloat `json:"slash_fraction_downtime"`
}

func MustGetParams(t testing.TB) Params {
	var params Params
	params.DowntimeJailDuration.UnmarshalText(nil)
	cmd.MustQuery(t, &params, "slashing", "params")
	require.Greater(t, params.DowntimeJailDuration.StdDuration(), 0*time.Second)
	require.False(t, params.MinSignedPerWindow.IsZero())
	require.Greater(t, params.SignedBlocksWindow.Int64(), int64(0))
	require.False(t, params.SlashFractionDoubleSign.IsZero())
	require.False(t, params.SlashFractionDowntime.IsZero())
	return params
}
