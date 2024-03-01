package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
)

var HomeDir string

type Tx struct {
	Code   int    `json:"code"`
	Hash   string `json:"txhash"`
	RawLog string `json:"raw_log"`
}

func Init(homeDir string) error {
	HomeDir = homeDir
	return os.Setenv("PATH", os.Getenv("HOME")+"/go/bin:"+os.Getenv("PATH"))
}

func MustInit(t testutil.TestingT, homeDir string) {
	err := Init(homeDir)
	require.NoError(t, err)
}

type ExecError struct {
	err    error
	output []byte
}

func MakeExecError(err error, output []byte) ExecError {
	return ExecError{err, output}
}

func (err ExecError) Error() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "[ERR] %s\n", err.err)
	fmt.Fprintf(&sb, "[OUT] %s", string(err.output))
	return sb.String()
}

func (err ExecError) Unwrap() error {
	return err.err
}

func (err ExecError) Output() []byte {
	return err.output
}

func Exec(name string, args ...string) ([]byte, error) {
	args = append(args, "--home="+HomeDir)
	cmd := exec.Command(name, args...)
	fmt.Println("[CMD]", cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, MakeExecError(err, output)
	}
	return output, nil
}

func MustExec(t testutil.TestingT, name string, args ...string) []byte {
	output, err := Exec(name, args...)
	require.NoError(t, err)
	return output
}

func Query(v any, args ...string) error {
	args = append([]string{"query"}, args...)
	args = append(args, "--output=json")
	output, err := Exec("titand", args...)
	if err != nil {
		return err
	}
	return json.Unmarshal(output, v)
}

func MustQuery(t testutil.TestingT, v any, args ...string) {
	err := Query(v, args...)
	require.NoError(t, err)
}

// Scan for the first line that contains JSON object and unmarshal
func UnmarshalJSON(data []byte, v any) error {
	s := bufio.NewScanner(bytes.NewBuffer(data))
	for s.Scan() {
		b := s.Bytes()
		if (b[0] == '[' && b[len(b)-1] == ']') || (b[0] == '{' && b[len(b)-1] == '}') {
			return json.Unmarshal(b, v)
		}
	}
	return fmt.Errorf("cannot unmarshal %s from: %s", reflect.TypeOf(v), string(data))
}

// Execute a command and write its output to w
func ExecWrite(w io.Writer, name string, args ...string) (*os.ProcessState, error) {
	cmd := exec.Command(name, args...)
	fmt.Fprintln(w, cmd.String())
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	type result struct {
		state *os.ProcessState
		err   error
	}
	done := make(chan result)
	go func() {
		state, err := cmd.Process.Wait()
		done <- result{state: state, err: err}
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
	r := <-done
	return r.state, r.err
}

// Must execute a command and write its output to w, panics on error
func MustExecWrite(t testutil.TestingT, w io.Writer, name string, args ...string) {
	state, err := ExecWrite(w, name, args...)
	require.NoError(t, err)
	require.Equal(t, 0, state.ExitCode(), state.String())
}
