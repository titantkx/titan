package testutil

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type PublicKey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

func (pk PublicKey) String() string {
	txt, err := json.Marshal(pk)
	if err != nil {
		panic(err)
	}
	return string(txt)
}

func MustGenerateEd25519PK(t testing.TB) PublicKey {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	require.NoError(t, err)
	base64.StdEncoding.EncodeToString(b)
	pk := PublicKey{
		Type: "/cosmos.crypto.ed25519.PubKey",
		Key:  base64.StdEncoding.EncodeToString(b),
	}
	return pk
}
