package status

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
)

type Status struct {
	NodeInfo      NodeInfo      `json:"NodeInfo"`
	SyncInfo      SyncInfo      `json:"SyncInfo"`
	ValidatorInfo ValidatorInfo `json:"ValidatorInfo"`
}

type NodeInfo struct{}

type SyncInfo struct {
	LastestBlockHash    string       `json:"latest_block_hash"`
	LatestAppHash       string       `json:"latest_app_hash"`
	LatestBlockHeight   testutil.Int `json:"latest_block_height"`
	LatestBlockTime     time.Time    `json:"latest_block_time"`
	EarliestBlockHash   string       `json:"earliest_block_hash"`
	EarliestAppHash     string       `json:"earliest_app_hash"`
	EarliestBlockHeight testutil.Int `json:"earliest_block_height"`
	EarliestBlockTime   time.Time    `json:"earliest_block_time"`
	CatchingUp          bool         `json:"catching_up"`
}

type ValidatorInfo struct{}

func MustGetStatus(t testing.TB) Status {
	output := cmd.MustExec(t, "titand", "status")
	var status Status
	err := json.Unmarshal(output, &status)
	require.NoError(t, err)
	return status
}

func MustGetLatestBlockHeight(t testing.TB) int64 {
	return MustGetStatus(t).SyncInfo.LatestBlockHeight.Int64()
}

// Wait until block height
func MustWait(t testing.TB, height int64) {
	for {
		curHeight := MustGetLatestBlockHeight(t)
		if curHeight >= height {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
