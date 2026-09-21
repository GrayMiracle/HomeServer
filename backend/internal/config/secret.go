// Secret key for JWT

package config

import (
	"crypto/rand" // crypto random for secure key generation
	"encoding/base64" // base64 encoding for storing key as string
	"os" // os for file handling
)

const secretKeyPath = "secret.key"

// Function to return server JWT signing secret key as bytes. First time generates new random secret, afterward reads from file so existing tokens validate as usual.
func HandleSecret() ([]byte, error) {

	// Try existing file if exists
	data, err := os.ReadFile(secretKeyPath)
	if err == nil {
		decoded, decodeErr := base64.StdEncoding.DecodeString(string(data))
		if decodeErr != nil {
			return nil, decodeErr
		}
		return decoded, nil
	}

	// If file doesn't exist, create new 32 byte key
	rawSecret := make([]byte, 32)
	if _, err := rand.Read(rawSecret); err != nil {
		return nil, err
	}

	// Encode raw secret as base64 string while storing
	encoded := base64.StdEncoding.EncodeToString(rawSecret)
	if err := os.WriteFile(secretKeyPath, []byte(encoded), 0600); err != nil {
		return nil, err
	}

	return rawSecret, nil
}