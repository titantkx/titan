package cmd_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/utils"
)

func TestMain(m *testing.M) {
	utils.InitSDKConfig()

	appPath, err := filepath.Abs("../../..")
	if err != nil {
		panic(err)
	}
	homePath, err := filepath.Abs(".titand")
	if err != nil {
		panic(err)
	}
	configPath, err := filepath.Abs("config.yml")
	if err != nil {
		panic(err)
	}

	if err := cmd.Init(homePath); err != nil {
		panic(err)
	}

	logger, err := os.Create("titand.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	fmt.Println("Blockchain starting...")
	process := startBlockchain(logger, appPath, homePath, configPath)
	fmt.Println("Blockchain started")

	done := make(chan struct{})

	go func() {
		state, err := process.Wait()
		if err != nil {
			panic(err)
		}
		if state.ExitCode() != 0 {
			panic(state.String())
		}
		done <- struct{}{}
	}()

	code := m.Run()

	process.Signal(os.Interrupt)
	<-done

	os.Exit(code)
}

func startBlockchain(w io.Writer, appPath, homePath, configPath string) *os.Process {
	cmd := exec.Command("ignite", "chain", "serve", "--skip-proto", "--reset-once", "--path="+appPath, "--home="+homePath, "--config="+configPath, "-v")
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
		isRunning := false
		r := bufio.NewReader(output)
		for {
			line, isPrefix, err := r.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			w.Write(line)
			if !isPrefix {
				w.Write([]byte("\n"))
			}
			if !isRunning && strings.Contains(string(line), "executed block") {
				isRunning = true
				ready <- struct{}{}
			}
		}
	}()
	<-ready
	return cmd.Process
}
