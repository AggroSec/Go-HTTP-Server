package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	refresh_token := make([]byte, 32)
	rand.Read(refresh_token)
	return hex.EncodeToString(refresh_token)
}
