package cmd_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/utils"
)

func TestMain(m *testing.M) {
	utils.InitSDKConfig()

	if err := os.MkdirAll("tmp", os.ModePerm); err != nil {
		panic(err)
	}

	rootDir, err := filepath.Abs("../../..")
	if err != nil {
		panic(err)
	}
	homePath, err := filepath.Abs("tmp/val1/.titand")
	if err != nil {
		panic(err)
	}

	if err := cmd.Init(homePath); err != nil {
		panic(err)
	}

	logger, err := os.Create("tmp/titand.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	stopChain(logger) // Stop any running instance

	fmt.Println("Installing titand...")
	install(logger, rootDir)

	fmt.Println("Building image...")
	buildImage(logger, rootDir)

	fmt.Println("Initializing blockchain...")
	initChain(logger)

	fmt.Println("Starting blockchain...")
	ready, done := startChain(logger)

	select {
	case <-ready:
		fmt.Println("Started blockchain")
	case <-done:
		panic("Blockchain is stopped before ready")
	}

	code := m.Run()

	stopChain(logger)

	<-done
	fmt.Println("Stopped blockchain")

	os.Exit(code)
}

func install(w io.Writer, rootDir string) {
	cwd := cmd.Getwd()
	cmd.Chdir(rootDir)
	defer cmd.Chdir(cwd)
	cmd.MustExecWrite(w, "make", "build")
	cmd.MustExecWrite(w, "cp", rootDir+"/build/titand", cmd.UserHomeDir()+"/go/bin")
}

func buildImage(w io.Writer, rootDir string) {
	cmd.MustExecWrite(w, "docker", "build", rootDir, "-t=e2e/titand")
}

func initChain(w io.Writer) {
	cmd.MustExecWrite(w, "sh", "init.sh")
}

func startChain(w io.Writer) (ready <-chan struct{}, done <-chan struct{}) {
	readyCh := make(chan struct{})
	doneCh := make(chan struct{})

	s := testutil.NewStreamer()

	go func() {
		cmd.MustExecWrite(s, "docker", "compose", "up")
		s.Close()
		doneCh <- struct{}{}
	}()

	go func() {
		isRunning := false
		r := bufio.NewReader(s)
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
				fmt.Fprintln(w)
			}
			if !isRunning && strings.Contains(string(line), "executed block") {
				isRunning = true
				readyCh <- struct{}{}
			}
		}
	}()

	return readyCh, doneCh
}

func stopChain(w io.Writer) {
	cmd.MustExecWrite(w, "docker", "compose", "down")
}
