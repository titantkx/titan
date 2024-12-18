package sample

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	mathrand "math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-faker/faker/v4"
)

// AccAddress returns a sample account address
func AccAddress() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr)
}

func ClassId() string {
	return strconv.FormatUint(mathrand.Uint64(), 10)
}

func Name() string {
	return faker.Name()
}

func Word() string {
	return faker.Word()
}

func Sentence() string {
	return faker.Sentence()
}

func Paragraph() string {
	return faker.Paragraph()
}

func URL() string {
	return faker.URL()
}

func Hash() string {
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func JSON() string {
	n := mathrand.Intn(10)
	m := make(map[string]interface{})
	for i := 0; i < n; i++ {
		m[Word()] = Sentence()
	}
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(b)
}
