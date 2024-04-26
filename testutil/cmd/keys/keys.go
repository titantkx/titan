package keys

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
)

type Key struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Address  string    `json:"address"`
	PubKey   PublicKey `json:"pubkey"`
	Mnemonic string    `json:"mnemonic"`
}

type PublicKey testutil.SinglePublicKey

func (pk *PublicKey) UnmarshalText(txt []byte) error {
	if len(txt) > 2 && txt[0] == '"' && txt[len(txt)-1] == '"' {
		txt = txt[1 : len(txt)-1]
	}
	return json.Unmarshal(txt, (*testutil.SinglePublicKey)(pk))
}

func MustShowAddress(t testutil.TestingT, name string) string {
	output := cmd.MustExec(t, "titand", "keys", "show", name, "--address")
	return strings.TrimSuffix(string(output), "\n")
}

func MustShow(t testutil.TestingT, name string) Key {
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

func MustAdd(t testutil.TestingT, name string) Key {
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

func MustDelete(t testutil.TestingT, name string) {
	cmd.MustExec(t, "titand", "keys", "delete", name, "-y")
}

func MustRename(t testutil.TestingT, oldName string, newName string) {
	cmd.MustExec(t, "titand", "keys", "rename", oldName, newName, "-y")
}

func MustList(t testutil.TestingT) []Key {
	output := cmd.MustExec(t, "titand", "keys", "list", "--output=json")
	require.NotNil(t, output)
	var keys []Key
	err := cmd.UnmarshalJSON(output, &keys)
	require.NoError(t, err)
	return keys
}

func MustExport(t testutil.TestingT, name string, password string) []byte {
	//nolint:gosec // this for testing purpose only
	command := exec.Command("titand", "keys", "export", name, "--home="+cmd.HomeDir, "--keyring-backend=test")
	stdin, err := command.StdinPipe()
	if err != nil {
		require.NoError(t, err)
	}
	go func() {
		defer stdin.Close()
		//nolint:errcheck
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

func MustImport(t testutil.TestingT, name string, fileName string, password string) {
	//nolint:gosec // this for testing purpose only
	command := exec.Command("titand", "keys", "import", name, fileName, "--home="+cmd.HomeDir, "--keyring-backend=test")
	stdin, err := command.StdinPipe()
	if err != nil {
		require.NoError(t, err)
	}
	go func() {
		defer stdin.Close()
		//nolint:errcheck
		io.WriteString(stdin, password)
	}()
	fmt.Println("[CMD]", command)
	output, err := command.CombinedOutput()
	if err != nil {
		err = cmd.MakeExecError(err, output)
	}
	require.NoError(t, err)
}

type MultisigKey struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Address string            `json:"address"`
	PubKey  MultisigPublicKey `json:"pubkey"`
}

type MultisigPublicKey testutil.MultisigPublicKey

func (pk *MultisigPublicKey) UnmarshalText(txt []byte) error {
	if len(txt) > 2 && txt[0] == '"' && txt[len(txt)-1] == '"' {
		txt = txt[1 : len(txt)-1]
	}
	return json.Unmarshal(txt, (*testutil.MultisigPublicKey)(pk))
}

func MustAddMultisig(t testutil.TestingT, name string, threshold int, keys ...string) MultisigKey {
	output := cmd.MustExec(t, "titand", "keys", "add", name, "--multisig="+strings.Join(keys, ","), "--multisig-threshold="+strconv.Itoa(threshold), "--output=json")
	require.NotNil(t, output)
	var mutisigKey MultisigKey
	err := cmd.UnmarshalJSON(output, &mutisigKey)
	require.NoError(t, err)
	require.Equal(t, name, mutisigKey.Name)
	require.Equal(t, "multi", mutisigKey.Type)
	require.NotEmpty(t, mutisigKey.Address)
	require.Equal(t, "/cosmos.crypto.multisig.LegacyAminoPubKey", mutisigKey.PubKey.Type)
	require.Equal(t, threshold, mutisigKey.PubKey.Threshold)
	require.Len(t, mutisigKey.PubKey.PublicKeys, len(keys))
	for i := range keys {
		expectedPublicKey := testutil.SinglePublicKey(MustShow(t, keys[i]).PubKey)
		require.Contains(t, mutisigKey.PubKey.PublicKeys, expectedPublicKey)
	}
	return mutisigKey
}
