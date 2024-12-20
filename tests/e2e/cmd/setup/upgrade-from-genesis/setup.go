package upgradefromgenesis

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/titantkx/titan/tests/e2e/cmd/setup"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
	"github.com/titantkx/titan/testutil/cmd/keys"
)

const (
	UpgradeName     = "v3_0_0"
	GenesisFileName = "genesis_mainnet_zero.json"
)

func Setup(m *testing.M, rootDir string, logger io.Writer) {
	t := testutil.NewMockTest(os.Stderr)
	defer t.Finish()
	testutil.HandleOSInterrupt(func() {
		setup.StopChain(t, logger, "docker-compose-genesis.yml")
		setup.StopChain(t, logger, "docker-compose-local.yml")
	})

	testutil.Chdir(t, "setup/upgrade-from-genesis")
	testutil.MkdirAll(t, "tmp", os.ModePerm)
	homeDir := testutil.AbsPath(t, "tmp/val1/.titand")
	cmd.MustInit(t, homeDir)

	// Check if genesis.json exists
	f, err := os.Open(GenesisFileName)
	if err != nil {
		panic("Cannot open genesis.json: " + err.Error())
	}
	f.Close()

	fmt.Println("Installing titand...")
	setup.Install(t, logger, rootDir)

	fmt.Println("Building image...")
	setup.BuildImage(t, logger, rootDir, "local")

	setup.StopChain(t, logger, "docker-compose-genesis.yml") // Stop any running instance

	fmt.Println("Initializing blockchain...")
	cmd.MustExecWrite(t, logger, "sh", "init.sh", GenesisFileName)

	fmt.Println("Starting blockchain...")
	ready, upgrade, done := setup.StartChainAndListenForUpgrade(t, logger, "docker-compose-genesis.yml", UpgradeName)

	select {
	case <-ready:
		fmt.Println("Started blockchain")
	case <-done:
		panic("Blockchain is stopped before ready")
	}

	fmt.Println("Upgrading blockchain...")
	val1 := keys.MustShowAddress(t, "val1")
	val2 := keys.MustShowAddress(t, "val2")
	setup.UpgradeChain(t, UpgradeName, val1, val2)

	<-upgrade
	fmt.Println("Ready to upgrade blockchain")

	fmt.Println("Restarting blockchain...")
	setup.StopChain(t, logger, "docker-compose-genesis.yml")
	ready, done = setup.StartChain(t, logger, "docker-compose-local.yml")

	select {
	case <-ready:
		fmt.Println("Restarted blockchain")
	case <-done:
		panic("Blockchain is stopped before ready")
	}

	code := m.Run()

	setup.StopChain(t, logger, "docker-compose-local.yml")

	//nolint:gocritic // We need to exit with the code
	os.Exit(code)
}
