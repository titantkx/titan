package block

import (
	"encoding/json"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
)

type Block struct {
	Header BlockHeader `json:"header"`
}

//nolint:revive
type BlockHeader struct {
	Time time.Time `json:"time"`
}

func MustGetBlockTime(t testutil.TestingT, height int64) time.Time {
	output := cmd.MustExec(t, "titand", "query", "block", testutil.FormatInt(height))
	var v struct {
		Block Block `json:"block"`
	}
	err := json.Unmarshal(output, &v)
	require.NoError(t, err)
	return v.Block.Header.Time
}
