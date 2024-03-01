package testutil

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/require"
)

func AbsPath(t TestingT, path string) string {
	absPath, err := filepath.Abs(path)
	require.NoError(t, err)
	return absPath
}

func MkdirAll(t TestingT, path string, perm fs.FileMode) {
	err := os.MkdirAll(path, perm)
	require.NoError(t, err)
}

func Chdir(t TestingT, dir string) {
	err := os.Chdir(dir)
	require.NoError(t, err)
}

func Getwd(t TestingT) string {
	wd, err := os.Getwd()
	require.NoError(t, err)
	return wd
}

func UserHomeDir(t TestingT) string {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	return homeDir
}
