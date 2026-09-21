// JWT handling

// auth package for authentication files
package auth

import (
	"time" // time for token expiration
	"github.com/golang-jwt/jwt/v5" // jwt handling
)

// Access token
const AccessTokenLifetime = 15 * time.Minute

// Refresh token
const RefreshTokenLifetime = 30 * 24 * time.Hour

// Claims, the data in the access token that ISNT encrypted
type Claims struct {
	UserID int `json:"user_id"`
	DeviceID string `json:"device_id"` // DeviceID for session management and logging
	jwt.RegisteredClaims
}

// GenerateAccessToken to create the signed JWT for a user
func GenerateAccessToken(userID int, deviceID string, secret []byte) (string, error) {
	claims := Claims{
		UserID: userID,
		DeviceID: deviceID, // Set the DeviceID when generating the token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ValidateAccessToken to check if the token is valid and return the claims
func ValidateAccessToken(tokenString string, secret []byte) (*Claims, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// Defend against algorithm-confusion attacks: explicitly check, the algorithm matches what we expect before trusting the secret.
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}
	return claims, nil
}