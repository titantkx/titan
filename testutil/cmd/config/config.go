package config

import (
	"encoding/json"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
)

type Config struct {
	ChainID        string `json:"chain-id"`
	KeyringBackend string `json:"keyring-backend"`
	Output         string `json:"output"`
	Node           string `json:"node"`
	BroadcastMode  string `json:"broadcast-mode"`
}

func MustGetConfig(t testutil.TestingT) Config {
	output := cmd.MustExec(t, "titand", "config")
	var config Config
	err := json.Unmarshal(output, &config)
	require.NoError(t, err)
	return config
}
