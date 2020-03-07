package rand

import (
	"crypto/rand"
	"encoding/base64"
)

// RememberTokenBytes represents the byte count
// for the generated string from String()
const RememberTokenBytes = 32

// Bytes generates n random bytes or an error
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String generates a byte slice of nBytes
// and returns a base64 URL encoded version
// of the slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken generates remember tokens
// of a predetermined byte size
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
