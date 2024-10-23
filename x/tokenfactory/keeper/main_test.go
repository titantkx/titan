package keeper_test

import (
	"os"
	"testing"

	"github.com/titantkx/titan/utils"
)

func TestMain(m *testing.M) {
	utils.InitSDKConfig()
	code := m.Run()
	os.Exit(code)
}
