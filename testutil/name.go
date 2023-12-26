package testutil

import (
	"crypto/rand"
	"encoding/hex"
)

func GetName() string {
	b := make([]byte, 16) // Increase the size to make the string longer
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
