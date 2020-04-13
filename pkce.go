package main

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"time"
)

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGH" +
	"IJKLMNOPQRSTUVWXYZ0123456789-._~"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func createCodeVerifier(length int) string {
	codeVerifier := generateRandomString(length)
	return codeVerifier
}

func createCodeChallenge(codeVerifier string, codeChallengeMethod string) string {
	if codeChallengeMethod == "plain" {
		return codeVerifier
	}
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
