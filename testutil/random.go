package testutil

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/stretchr/testify/require"
)

func MustRandomString(t TestingT, n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return base64.StdEncoding.EncodeToString(b)[:n]
}
