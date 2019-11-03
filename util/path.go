package util

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
)

var (
	// ErrorHashingFailed is returned if the hashing of the provided path fails
	ErrorHashingFailed = errors.New("hhashing the provided path failed")
)

// z-base-32 charset only, easier to remember and type
const pathChars = "ybndrfg8ejkmcpqxotluwisza345h769"

// GenerateZBase32RandomPath generates a file path for Human ease using z-base-32
func GenerateZBase32RandomPath(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = pathChars[rand.Intn(len(pathChars))]
	}

	return string(b)
}

// HashForPath generates a hash used to store the file so that the user-generated path is not stored or used
func HashForPath(path string) (string, error) {
	new256 := sha256.New()
	_, err := new256.Write([]byte(path))

	if err != nil {
		fmt.Println(err)
		return "", ErrorHashingFailed
	}

	out := new256.Sum(nil)

	return fmt.Sprintf("%x", out), nil
}
