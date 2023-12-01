package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
		return output, ExecError{err, output}
	}
	return output, nil
}

func MustExec(t testing.TB, name string, args ...string) []byte {
	output, err := Exec(name, args...)
	if err != nil {
		t.Fatal(err)
	}
	return output
}

func MustQuery(t testing.TB, v any, args ...string) {
	args = append([]string{"query"}, args...)
	args = append(args, "--output=json")
	output := MustExec(t, "titand", args...)
	err := json.Unmarshal(output, v)
	require.NoError(t, err)
}
