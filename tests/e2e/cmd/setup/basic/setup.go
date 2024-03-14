package basic

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/titantkx/titan/tests/e2e/cmd/setup"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
)

func Setup(m *testing.M, rootDir string, logger io.Writer) {
	t := testutil.NewMockTest(os.Stderr)
	defer t.Finish()

	testutil.Chdir(t, "setup/basic")
	testutil.MkdirAll(t, "tmp", os.ModePerm)
	homeDir := testutil.AbsPath(t, "tmp/val1/.titand")
	cmd.MustInit(t, homeDir)

	setup.StopChain(t, logger, "docker-compose.yml") // Stop any running instance

	fmt.Println("Installing titand...")
	setup.Install(t, logger, rootDir)

	fmt.Println("Building image...")
	setup.BuildImage(t, logger, rootDir, "latest")

	fmt.Println("Initializing blockchain...")
	cmd.MustExecWrite(t, logger, "sh", "init.sh")

	fmt.Println("Starting blockchain...")
	ready, done := setup.StartChain(t, logger, "docker-compose.yml")

	select {
	case <-ready:
		fmt.Println("Started blockchain")
	case <-done:
		panic("Blockchain is stopped before ready")
	}

	code := m.Run()

	setup.StopChain(t, logger, "docker-compose.yml")

	<-done
	fmt.Println("Stopped blockchain")

	os.Exit(code)
}
