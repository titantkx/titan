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
	cwd := Getwd()
	Chdir(w, rootDir)
	defer Chdir(w, cwd)
	Exec(w, "make", "build")
	Exec(w, "cp", rootDir+"/build/titand", UserHomeDir()+"/go/bin")
}

func buildImage(w io.Writer, rootDir string) {
	Exec(w, "docker", "build", rootDir, "-t=e2e/titand")
}

func initChain(w io.Writer) {
	Exec(w, "sh", "init.sh")
}

func startChain(w io.Writer) (ready <-chan struct{}, done <-chan struct{}) {
	cmd := exec.Command("docker", "compose", "up")
	fmt.Fprintln(w, cmd.String())
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	readyCh := make(chan struct{})
	doneCh := make(chan struct{})
	go func() {
		state, err := cmd.Process.Wait()
		if err != nil {
			panic(err)
		}
		if state.ExitCode() != 0 {
			panic(state.String())
		}
		doneCh <- struct{}{}
	}()
	go func() {
		isRunning := false
		r := bufio.NewReader(stdoutPipe)
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
	go func() {
		r := bufio.NewReader(stderrPipe)
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
		}
	}()
	return readyCh, doneCh
}

func stopChain(w io.Writer) {
	Exec(w, "docker", "compose", "down")
}

func Exec(w io.Writer, name string, args ...string) {
	cmd := exec.Command(name, args...)
	fmt.Fprintln(w, cmd.String())
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	done := make(chan struct{})
	go func() {
		state, err := cmd.Process.Wait()
		if err != nil {
			panic(err)
		}
		if state.ExitCode() != 0 {
			panic(state.String())
		}
		done <- struct{}{}
	}()
	go func() {
		r := bufio.NewReader(stdoutPipe)
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
		}
	}()
	go func() {
		r := bufio.NewReader(stderrPipe)
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
		}
	}()
	<-done
}

func Chdir(w io.Writer, dir string) {
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "cd %s\n", dir)
}

func Getwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

func UserHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}
