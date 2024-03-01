package testutil

import (
	"os"

	"github.com/stretchr/testify/require"
)

func MustCreateTemp(t TestingT, name string) *os.File {
	file, err := os.CreateTemp("", name)
	require.NoError(t, err)
	require.NotNil(t, file)
	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})
	return file
}
