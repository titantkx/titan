package keys

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil/cmd"
)

func MustGetAddress(t testing.TB, name string) string {
	addr := cmd.MustExec(t, "titand", "keys", "show", name, "--address")
	require.NotEmpty(t, addr)
	return strings.TrimSuffix(string(addr), "\n")
}

func MustCreateAccount(t testing.TB, name string) string {
	cmd.MustExec(t, "titand", "keys", "add", name)
	return MustGetAddress(t, name)
}

func MustDeleteAccount(t testing.TB, name string) {
	cmd.MustExec(t, "titand", "keys", "delete", name, "-y")
}
