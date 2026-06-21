package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	encodedString := hex.EncodeToString(randomBytes)
	return encodedString
}
