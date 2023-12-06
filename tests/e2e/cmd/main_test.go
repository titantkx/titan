package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tokenize-titan/titan/testutil/cmd"
)

func TestMain(m *testing.M) {
	homeDir, err := filepath.Abs("../../../local_test_data/.titan_val1")
	if err != nil {
		panic(err)
	}
	if err := cmd.Init(homeDir); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
