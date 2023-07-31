package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// TODO 可以利用库来做加密,这样就不用自己来加盐了(比如用go-password-encoder)

func GetSHA256HashCode(message string) string {
	bytes := sha256.Sum256([]byte(message))
	hashCode := hex.EncodeToString(bytes[:])
	return hashCode
}
