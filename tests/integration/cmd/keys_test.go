package cmd_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/testutil/cmd/keys"
)

func MustAddKey(t testing.TB) keys.Key {
	name := testutil.GetName()
	t.Cleanup(func() {
		defer keys.MustDelete(t, name)
		testutil.PutName(name)
	})
	return keys.MustAdd(t, name)
}

func TestAddKey(t *testing.T) {
	t.Parallel()

	MustAddKey(t)
}

func TestAddKeyDuplicatedName(t *testing.T) {
	t.Parallel()

	key := MustAddKey(t)

	_, err := cmd.Exec("titand", "keys", "add", key.Name)

	require.Error(t, err)
	require.ErrorContains(t, err, "Error: EOF")
}

func TestShowKey(t *testing.T) {
	t.Parallel()

	key := MustAddKey(t)

	keys.MustShow(t, key.Name)
	keys.MustShow(t, key.Address)
}

func TestShowKeyAddress(t *testing.T) {
	t.Parallel()

	key := MustAddKey(t)

	address := keys.MustShowAddress(t, key.Name)

	require.Equal(t, key.Address, address)
}

func TestShowKeyNotFound(t *testing.T) {
	t.Parallel()

	name := testutil.GetName()
	defer testutil.PutName(name)

	_, err := cmd.Exec("titand", "keys", "show", name)

	require.Error(t, err)
	require.ErrorContains(t, err, "not a valid name or address")
}

func TestDeleteKey(t *testing.T) {
	t.Parallel()

	name := testutil.GetName()
	defer testutil.PutName(name)

	keys.MustAdd(t, name)
	keys.MustDelete(t, name)

	_, err := cmd.Exec("titand", "keys", "show", name)

	require.Error(t, err)
	require.ErrorContains(t, err, "not a valid name or address")
}

func TestDeleteKeyNotFound(t *testing.T) {
	t.Parallel()

	name := testutil.GetName()
	defer testutil.PutName(name)

	_, err := cmd.Exec("titand", "keys", "delete", name)

	require.Error(t, err)
	require.ErrorContains(t, err, "key not found")
}

func TestRenameKey(t *testing.T) {
	t.Parallel()

	oldName := testutil.GetName()
	defer testutil.PutName(oldName)
	newName := testutil.GetName()
	defer testutil.PutName(newName)

	defer keys.MustDelete(t, newName)
	oldKey := keys.MustAdd(t, oldName)
	keys.MustRename(t, oldName, newName)
	newKey := keys.MustShow(t, newName)

	require.Equal(t, oldKey.Type, newKey.Type)
	require.Equal(t, oldKey.Address, newKey.Address)
	require.Equal(t, oldKey.PubKey.Type, newKey.PubKey.Type)
	require.Equal(t, oldKey.PubKey.Key, newKey.PubKey.Key)
}

func TestRenameKeyNotFound(t *testing.T) {
	t.Parallel()

	oldName := testutil.GetName()
	defer testutil.PutName(oldName)
	newName := testutil.GetName()
	defer testutil.PutName(newName)

	_, err := cmd.Exec("titand", "keys", "rename", oldName, newName)

	require.Error(t, err)
	require.ErrorContains(t, err, "key not found")
}

func TestRenameKeyToExistingKey(t *testing.T) {
	t.Parallel()

	key1 := MustAddKey(t)
	key2 := MustAddKey(t)

	_, err := cmd.Exec("titand", "keys", "rename", key1.Name, key2.Name)

	require.Error(t, err)
	require.ErrorContains(t, err, "Error: EOF")
}

func TestListKeys(t *testing.T) {
	t.Parallel()

	expectedKey := MustAddKey(t)

	keyList := keys.MustList(t)

	var actualKey keys.Key
	for _, key := range keyList {
		if key.Name == expectedKey.Name {
			actualKey = key
		}
	}

	require.Equal(t, expectedKey.Name, actualKey.Name)
	require.Equal(t, expectedKey.Type, actualKey.Type)
	require.Equal(t, expectedKey.Address, actualKey.Address)
	require.Equal(t, expectedKey.PubKey.Type, actualKey.PubKey.Type)
	require.Equal(t, expectedKey.PubKey.Key, actualKey.PubKey.Key)
}

func exportKey(t testing.TB, password string, w io.Writer) keys.Key {
	name := testutil.GetName()
	defer testutil.PutName(name)
	defer keys.MustDelete(t, name)
	key := keys.MustAdd(t, name)
	output := keys.MustExport(t, name, password)
	_, err := w.Write(output)
	require.NoError(t, err)
	return key
}

func TestExportKey(t *testing.T) {
	t.Parallel()

	password := testutil.MustRandomString(t, 12)
	exportKey(t, password, io.Discard)
}

func TestImportKey(t *testing.T) {
	t.Parallel()

	password := testutil.MustRandomString(t, 12)
	file, err := os.CreateTemp("", "private_key_*.txt")
	require.NoError(t, err)
	require.NotNil(t, file)
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	exportedKey := exportKey(t, password, file)

	name := testutil.GetName()
	defer testutil.PutName(name)
	defer keys.MustDelete(t, name)
	keys.MustImport(t, name, file.Name(), password)

	importedKey := keys.MustShow(t, name)

	require.Equal(t, exportedKey.Type, importedKey.Type)
	require.Equal(t, exportedKey.Address, importedKey.Address)
	require.Equal(t, exportedKey.PubKey.Type, importedKey.PubKey.Type)
	require.Equal(t, exportedKey.PubKey.Key, importedKey.PubKey.Key)
}
