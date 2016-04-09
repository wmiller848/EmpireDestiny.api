package util

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func Random(keySize int) []byte {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return key
}

func Hex(d []byte) string {
	return hex.EncodeToString(d)
}

func Unhex(d string) []byte {
	byts, err := hex.DecodeString(d)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return byts
}

func Hash(d string) string {
	hasher := sha1.New()
	hasher.Write([]byte(d))
	sha1 := hasher.Sum(nil)
	return Hex(sha1)
}
