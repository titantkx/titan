package testutil

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type PublicKey struct {
	Type  string `json:"@type"`
	Value any    `json:"-"`
}

func (pk *PublicKey) UnmarshalJSON(data []byte) error {
	var k struct {
		Type string `json:"@type"`
	}
	if err := json.Unmarshal(data, &k); err != nil {
		return err
	}
	pk.Type = k.Type
	switch pk.Type {
	case "/cosmos.crypto.multisig.LegacyAminoPubKey":
		var v MultisigPublicKey
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		pk.Value = v
	default:
		var v SinglePublicKey
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		pk.Value = v
	}
	return nil
}

type SinglePublicKey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

func (pk SinglePublicKey) GetType() string {
	return pk.Type
}

func (pk SinglePublicKey) String() string {
	txt, err := json.Marshal(pk)
	if err != nil {
		panic(err)
	}
	return string(txt)
}

func MustGenerateEd25519PK(t testing.TB) SinglePublicKey {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	require.NoError(t, err)
	base64.StdEncoding.EncodeToString(b)
	pk := SinglePublicKey{
		Type: "/cosmos.crypto.ed25519.PubKey",
		Key:  base64.StdEncoding.EncodeToString(b),
	}
	return pk
}

type MultisigPublicKey struct {
	Type       string            `json:"@type"`
	Threshold  int               `json:"threshold"`
	PublicKeys []SinglePublicKey `json:"public_keys"`
}

func (pk MultisigPublicKey) GetType() string {
	return pk.Type
}
