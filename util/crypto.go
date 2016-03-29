package util

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hex(d []byte) string {
	return hex.EncodeToString(d)
}

func Hash(d string) string {
	hasher := sha1.New()
	hasher.Write([]byte(d))
	sha1 := hasher.Sum(nil)
	return Hex(sha1)
}
