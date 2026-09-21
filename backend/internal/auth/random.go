// Random generation with crypto/rand for refresh token

package auth

import (
	"crypto/rand" // crypto random for secure token generation
	"encoding/base64" // base64 encoding for storing token as string
)

// GenerateRandomString for refreshtoken string
func generateRandomString(numBytes int) (string, error) {
	raw := make([]byte, numBytes)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(raw), nil
}

// GenerateRefreshToken with generateRandomString thats not a JWT, used to look up in database and verify against stored hash
func GenerateRefreshToken() (string, error) {
	return generateRandomString(48)
}