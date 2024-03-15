package cmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/titantkx/titan/tests/e2e/cmd/setup/basic"
	"github.com/titantkx/titan/tests/e2e/cmd/setup/upgrade"
	upgradefromgenesis "github.com/titantkx/titan/tests/e2e/cmd/setup/upgrade-from-genesis"
	"github.com/titantkx/titan/utils"
)

const (
	TestTypeBasic              = "basic"
	TestTypeUpgrade            = "upgrade"
	TestTypeUpgradeFromGenesis = "upgrade-from-genesis"
)

func TestMain(m *testing.M) {
	// defer os.Exit here to avoid hang when test fail in Setup step
	defer os.Exit(1)

	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}

	utils.InitSDKConfig()
	utils.RegisterDenoms()

	rootDir, err := filepath.Abs("../../..")
	if err != nil {
		panic(err)
	}

	// Always log to file except when LOG_OUTPUT_TYPE is set to "std"
	logOutputType := os.Getenv("LOG_OUTPUT_TYPE")
	logger := os.Stdout
	if logOutputType != "std" {
		logger, err = os.Create("titand.log")
		if err != nil {
			panic(err)
		}
	}

	defer logger.Close()

	testType := os.Getenv("TEST_TYPE")
	if testType == "" {
		testType = TestTypeBasic
	}

	switch testType {
	case TestTypeBasic:
		fmt.Printf("Test type: %s\n", testType)
		basic.Setup(m, rootDir, logger)
	case TestTypeUpgrade:
		fmt.Printf("Test type: %s\n", testType)
		upgrade.Setup(m, rootDir, logger)
	case TestTypeUpgradeFromGenesis:
		fmt.Printf("Test type: %s\n", testType)
		upgradefromgenesis.Setup(m, rootDir, logger)
	default:
		panic("Invalid test type: " + testType)
	}
}
