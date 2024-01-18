package cmd_test

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/utils"
)

func TestMain(m *testing.M) {
	utils.InitSDKConfig()
	appPath, err := filepath.Abs("../../..")
	if err != nil {
		panic(err)
	}
	homePath, err := filepath.Abs("./.titand")
	if err != nil {
		panic(err)
	}
	configPath, err := filepath.Abs("./config.yml")
	if err != nil {
		panic(err)
	}
	if err := cmd.Init(homePath); err != nil {
		panic(err)
	}
	process := startBlockchain(appPath, homePath, configPath)
	code := m.Run()
	process.Kill()
	os.Exit(code)
}

func startBlockchain(appPath, homePath, configPath string) *os.Process {
	cmd := exec.Command("ignite", "chain", "serve", "--skip-proto", "--reset-once", "--path="+appPath, "--home="+homePath, "--config="+configPath)
	fmt.Println("[CMD]", cmd)
	output, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	ready := make(chan struct{})
	go func() {
		state, err := cmd.Process.Wait()
		if err != nil {
			panic(err)
		}
		if state.ExitCode() != 0 {
			panic(state.String())
		}
	}()
	go func() {
		started := false
		scanner := bufio.NewScanner(output)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			if !started && strings.Contains(scanner.Text(), "Blockchain is running") {
				started = true
				ready <- struct{}{}
			}
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}()
	<-ready
	time.Sleep(3 * time.Second)
	return cmd.Process
}
