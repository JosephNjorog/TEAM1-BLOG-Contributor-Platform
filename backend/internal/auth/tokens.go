package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// generateOpaqueToken returns a URL-safe random token plus the sha256 hex
// digest that gets stored in the database (we never store the raw secret).
func generateOpaqueToken() (raw string, hash string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", err
	}
	raw = hex.EncodeToString(buf)
	hash = hashToken(raw)
	return raw, hash, nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
