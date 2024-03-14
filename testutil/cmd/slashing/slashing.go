package slashing

import (
	"time"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
)

type Params struct {
	DowntimeJailDuration    testutil.Duration `json:"downtime_jail_duration"`
	MinSignedPerWindow      testutil.Float    `json:"min_signed_per_window"`
	SignedBlocksWindow      testutil.Int      `json:"signed_blocks_window"`
	SlashFractionDoubleSign testutil.Float    `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   testutil.Float    `json:"slash_fraction_downtime"`
}

func MustGetParams(t testutil.TestingT) Params {
	var params Params
	cmd.MustQuery(t, &params, "slashing", "params")
	require.Greater(t, params.DowntimeJailDuration.StdDuration(), 0*time.Second)
	require.False(t, params.MinSignedPerWindow.IsZero())
	require.Greater(t, params.SignedBlocksWindow.Int64(), int64(0))
	require.False(t, params.SlashFractionDoubleSign.IsZero())
	require.False(t, params.SlashFractionDowntime.IsZero())
	return params
}
