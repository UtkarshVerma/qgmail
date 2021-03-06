package pkce

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomString generates a random string.
func RandomString(length int) string {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGH" +
		"IJKLMNOPQRSTUVWXYZ0123456789-._~"

	b := make([]byte, length)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

// CodeVerifier creates a code verifier for OAuth2.
func CodeVerifier(length int) string {
	return RandomString(length)
}

// CodeChallenge creates a code challenge for OAuth2.
func CodeChallenge(codeVerifier string, method string) string {
	if method == "plain" {
		return codeVerifier
	}

	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
