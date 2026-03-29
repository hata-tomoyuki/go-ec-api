package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

const refreshTokenBytes = 32

func newRefreshTokenPlaintext() (string, error) {
	b := make([]byte, refreshTokenBytes)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashRefreshToken(plaintext string) string {
	sum := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum[:])
}
