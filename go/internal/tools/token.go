package tools

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken creates a secure random token of n bytes (hex-encoded)
func GenerateToken(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
