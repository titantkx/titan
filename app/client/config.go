package client

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/titantkx/titan/app"
)

var (
	DefaultChainID = app.DefaultChainID
	DefaultHomeDir = app.DefaultNodeHome
)

const (
	DefaultKeyringBackend = keyring.BackendMemory
	DefaultOutputFormat   = "text"
	DefaultNodeURI        = "http://localhost:26657"
	DefaultBroadcastMode  = "sync"
)

type Config struct {
	ChainID        string
	HomeDir        string
	KeyringDir     string
	KeyringBackend string
	Output         string
	Node           string
	BroadcastMode  string
	GRPC           string
	GRPCInsecure   bool
}

func (conf Config) GetChainID() string {
	if conf.ChainID == "" {
		return DefaultChainID
	}
	return conf.ChainID
}

func (conf Config) GetHomeDir() string {
	if conf.HomeDir == "" {
		return DefaultHomeDir
	}
	return conf.HomeDir
}

func (conf Config) GetKeyringDir() string {
	if conf.KeyringDir == "" {
		return conf.GetHomeDir()
	}
	return conf.KeyringDir
}

func (conf Config) GetKeyringBackend() string {
	if conf.KeyringBackend == "" {
		return DefaultKeyringBackend
	}
	return conf.KeyringBackend
}

func (conf Config) GetOutputFormat() string {
	if conf.Output == "" {
		return DefaultOutputFormat
	}
	return conf.Output
}

func (conf Config) GetNodeURI() string {
	if conf.Node == "" {
		return DefaultNodeURI
	}
	return conf.Node
}

func (conf Config) GetBroadcastMode() string {
	if conf.BroadcastMode == "" {
		return DefaultBroadcastMode
	}
	return conf.BroadcastMode
}
