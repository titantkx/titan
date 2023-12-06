package keys

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
)

type Key struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Address  string    `json:"address"`
	PubKey   PublicKey `json:"pubkey"`
	Mnemonic string    `json:"mnemonic"`
}

type PublicKey testutil.PublicKey

func (pk *PublicKey) UnmarshalText(txt []byte) error {
	if len(txt) > 2 && txt[0] == '"' && txt[len(txt)-1] == '"' {
		txt = txt[1 : len(txt)-1]
	}
	return json.Unmarshal(txt, (*testutil.PublicKey)(pk))
}

func MustShowAddress(t testing.TB, name string) string {
	output := cmd.MustExec(t, "titand", "keys", "show", name, "--address")
	require.NotEmpty(t, output)
	return strings.TrimSuffix(string(output), "\n")
}

func MustShow(t testing.TB, name string) Key {
	output := cmd.MustExec(t, "titand", "keys", "show", name, "--output=json")
	require.NotNil(t, output)
	var key Key
	err := cmd.UnmarshalJSON(output, &key)
	require.NoError(t, err)
	require.Contains(t, []string{key.Name, key.Address}, name)
	require.Equal(t, "local", key.Type)
	require.NotEmpty(t, key.Address)
	require.Equal(t, "/ethermint.crypto.v1.ethsecp256k1.PubKey", key.PubKey.Type)
	require.NotEmpty(t, key.PubKey.Key)
	require.Empty(t, key.Mnemonic)
	return key
}

func MustAdd(t testing.TB, name string) Key {
	output := cmd.MustExec(t, "titand", "keys", "add", name, "--output=json")
	require.NotNil(t, output)
	var key Key
	err := cmd.UnmarshalJSON(output, &key)
	require.NoError(t, err)
	require.Equal(t, name, key.Name)
	require.Equal(t, "local", key.Type)
	require.NotEmpty(t, key.Address)
	require.Equal(t, "/ethermint.crypto.v1.ethsecp256k1.PubKey", key.PubKey.Type)
	require.NotEmpty(t, key.PubKey.Key)
	require.NotEmpty(t, key.Mnemonic)
	return key
}

func MustDelete(t testing.TB, name string) {
	cmd.MustExec(t, "titand", "keys", "delete", name, "-y")
}

func MustRename(t testing.TB, oldName string, newName string) {
	cmd.MustExec(t, "titand", "keys", "rename", oldName, newName, "-y")
}

func MustList(t testing.TB) []Key {
	output := cmd.MustExec(t, "titand", "keys", "list", "--output=json")
	require.NotNil(t, output)
	var keys []Key
	err := cmd.UnmarshalJSON(output, &keys)
	require.NoError(t, err)
	return keys
}

func MustExport(t testing.TB, name string, password string) []byte {
	command := exec.Command("titand", "keys", "export", name, "--home="+cmd.HomeDir, "--keyring-backend=test")
	stdin, err := command.StdinPipe()
	if err != nil {
		require.NoError(t, err)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, password)
	}()
	fmt.Println("[CMD]", command)
	output, err := command.CombinedOutput()
	if err != nil {
		err = cmd.MakeExecError(err, output)
	}
	require.NoError(t, err)
	require.NotEmpty(t, output)
	return output
}

func MustImport(t testing.TB, name string, fileName string, password string) {
	command := exec.Command("titand", "keys", "import", name, fileName, "--home="+cmd.HomeDir, "--keyring-backend=test")
	stdin, err := command.StdinPipe()
	if err != nil {
		require.NoError(t, err)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, password)
	}()
	fmt.Println("[CMD]", command)
	output, err := command.CombinedOutput()
	if err != nil {
		err = cmd.MakeExecError(err, output)
	}
	require.NoError(t, err)
}
