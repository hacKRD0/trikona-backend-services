package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomPassword generates a secure random password of specified length
func GenerateRandomPassword(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
} 