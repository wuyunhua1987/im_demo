package util

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandHex(len int) string {
	b := make([]byte, len)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
