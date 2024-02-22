package keeper_test

import (
	"os"
	"testing"

	"github.com/tokenize-titan/titan/utils"
)

func TestMain(m *testing.M) {
	utils.InitSDKConfig()
	code := m.Run()
	os.Exit(code)
}
