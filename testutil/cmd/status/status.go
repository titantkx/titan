package status

import (
	"encoding/json"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
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

func MustGetStatus(t testutil.TestingT) Status {
	output := cmd.MustExec(t, "titand", "status")
	var status Status
	err := json.Unmarshal(output, &status)
	require.NoError(t, err)
	return status
}

func MustGetLatestBlockHeight(t testutil.TestingT) int64 {
	return MustGetStatus(t).SyncInfo.LatestBlockHeight.Int64()
}

// Wait until block height
func MustWait(t testutil.TestingT, height int64) {
	for {
		curHeight := MustGetLatestBlockHeight(t)
		if curHeight >= height {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
