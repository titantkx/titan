package upgrade

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/tokenize-titan/titan/tests/e2e/cmd/setup"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/testutil/cmd/keys"
)

const UpgradeName = "v2_0_0"

func Setup(m *testing.M, rootDir string, logger io.Writer) {
	t := testutil.NewMockTest(os.Stderr)
	defer t.Finish()

	testutil.Chdir(t, "setup/upgrade")
	testutil.MkdirAll(t, "tmp", os.ModePerm)
	homeDir := testutil.AbsPath(t, "tmp/val1/.titand")
	cmd.MustInit(t, homeDir)

	fmt.Println("Installing titand...")
	setup.Install(t, logger, rootDir)

	fmt.Println("Building image...")
	setup.BuildImage(t, logger, rootDir, "latest")

	setup.StopChain(t, logger, "docker-compose-genesis.yml") // Stop any running instance

	fmt.Println("Initializing blockchain...")
	cmd.MustExecWrite(t, logger, "sh", "init.sh")

	fmt.Println("Starting blockchain...")
	ready, upgrade, done := setup.StartChainAndListenForUpgrade(t, logger, "docker-compose-genesis.yml", UpgradeName)

	select {
	case <-ready:
		fmt.Println("Started blockchain")
	case <-done:
		panic("Blockchain is stopped before ready")
	}

	fmt.Println("Seeding data...")
	seed(t)

	fmt.Println("Upgrading blockchain...")
	val1 := keys.MustShowAddress(t, "val1")
	val2 := keys.MustShowAddress(t, "val2")
	setup.UpgradeChain(t, UpgradeName, val1, val2)

	<-upgrade
	fmt.Println("Ready to upgrade blockchain")

	fmt.Println("Restarting blockchain...")
	setup.StopChain(t, logger, "docker-compose-genesis.yml")
	ready, done = setup.StartChain(t, logger, "docker-compose-latest.yml")

	select {
	case <-ready:
		fmt.Println("Restarted blockchain")
	case <-done:
		panic("Blockchain is stopped before ready")
	}

	code := m.Run()

	setup.StopChain(t, logger, "docker-compose-latest.yml")

	<-done
	fmt.Println("Stopped blockchain")

	os.Exit(code)
}
