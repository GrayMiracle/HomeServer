// Middleware for authentication

package auth

import (
	"net/http" // HTTP for middleware
	"fmt" // Format for Sprintf
)

// JWT Middleware that checks auth header for valid access token. If valid, allows request to real handler
func JWTMiddleware(jwtSecret []byte, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get access token from cookie
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// ValidateAccessToken to check validity
		claims, err := ValidateAccessToken(cookie.Value, jwtSecret)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// claims.UserID set in header so token doesnt need to be reparsed in order to find userid/who is making the request
		r.Header.Set("X-User-ID", fmt.Sprintf("%d", claims.UserID))

		r.Header.Set("X-Device-ID", claims.DeviceID) // Set device ID in header for logging and other uses

		// Finally call real handler
		next(w, r)
	}
}