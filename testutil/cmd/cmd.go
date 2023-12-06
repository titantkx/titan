package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
	args = append(args, "--home="+HomeDir, "--keyring-backend=test")
	cmd := exec.Command(name, args...)
	fmt.Println("[CMD]", cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, MakeExecError(err, output)
	}
	return output, nil
}

func MustExec(t testing.TB, name string, args ...string) []byte {
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

func MustQuery(t testing.TB, v any, args ...string) {
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
