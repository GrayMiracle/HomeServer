// Authentication request handler

package api

import (
	"encoding/json" // json handling
	"net/http"      // http handling
	"time"		  // time handling
	"strconv" // string conversion for userID in header
	"golang.org/x/crypto/bcrypt" // password hashing
	"homeserver/internal/db" // database handling
	"homeserver/internal/auth" // authentication handling
)

// Login request structure like the JSON body frontend sends
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DeviceID string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
}

// Login response sent back upon success
type loginResponse struct {
	RefreshToken string `json:"refreshToken"`
	UserID int `json:"userId"`
}

const isDevelopment = false // Set to false in production for secure cookies

// LoginHandler to handle login requests from POST /auth/login with jwtSecret passed from main.go
func LoginHandler(w http.ResponseWriter, r *http.Request, jwtSecret []byte) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed) // Only allow POST
		return
	}

	// Rate limit check
	if auth.IsRateLimited(r) {
		http.Error(w, "Too many failed attempts, try again later", http.StatusTooManyRequests)
		return
	}

	// new variable to unload JSON body into
	var req loginRequest
	// Fills req with json body here
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest) // If theres an error decoding JSON body throw it here
		return
	}

	// GetUserByUsername to fetch user, if doesnt exist throw error here
	userID, hashedPassword, err := db.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError) // Database error
		return
	}
	// If username is wrong
	if userID == 0 {
		auth.RecordFailedAttempt(r)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized) // User not found
		return
	}
	// If password is wrong
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		auth.RecordFailedAttempt(r)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized) // Wrong password
		return
	}

	// Password is correct, create access and refresh token
	accessToken, err := auth.GenerateAccessToken(userID, req.DeviceID, jwtSecret)
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError) // Error creating access token
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError) // Error creating refresh token
		return
	}

	// Store refresh token
	refreshHash, err:= bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError) // Error hashing refresh token
		return
	}

	expiresAt := time.Now().Add(auth.RefreshTokenLifetime)
	if err := db.CreateDevice(req.DeviceID, req.DeviceName, userID, string(refreshHash), expiresAt); err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	// Set access token as HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name: "access_token",
		Value: accessToken,
		HttpOnly: true,
		Secure: !isDevelopment,
		SameSite: http.SameSiteStrictMode,
		Path: "/",
		MaxAge: int(auth.AccessTokenLifetime.Seconds()),
	})

	// Successful login
	auth.RecordSuccessfulLogin(r)

	// Best effort logging, focusing on good logins over failed logs
	_ = db.LogActivity("login", "", strconv.Itoa(userID), req.DeviceID, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{
		RefreshToken: refreshToken,
		UserID:       userID,
	})
}

// refreshRequest similar to body sent to auth/refresh
type refreshRequest struct {
	UserID int `json:"userId"`
	DeviceID string `json:"deviceId"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshHandler for new access tokens
func RefreshHandler(w http.ResponseWriter, r *http.Request, jwtSecret []byte) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed) // Only allow POST
		return
	}

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest) // Error decoding JSON body
		return
	}

	device, err := db.VerifyRefreshToken(req.UserID, req.DeviceID, req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized) // Invalid refresh token
		return
	}
	if device == nil {
		http.Error(w, "Please login again", http.StatusUnauthorized) // Device not found or token invalid
		return
	}

	accessToken, err := auth.GenerateAccessToken(req.UserID, req.DeviceID, jwtSecret)
	if err != nil {
		http.Error(w, "Could not create access token", http.StatusInternalServerError) // Error creating access token
		return
	}

	// Set fresh cookie
	http.SetCookie(w, &http.Cookie {
		Name: "access_token",
		Value: accessToken,
		HttpOnly: true,
		Secure: !isDevelopment,
		SameSite: http.SameSiteStrictMode,
		Path: "/",
		MaxAge: int(auth.AccessTokenLifetime.Seconds()),
	})

	// Best effort timestamp update even in refresh failures
	_ = db.UpdateLastSeen(req.DeviceID)

	// Best effort logging, log even in refresh failures
	_ = db.LogActivity("refresh", "", strconv.Itoa(req.UserID), req.DeviceID, r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
}

// Logout struct similar to body sent to /auth/logout
type logoutRequest struct {
	UserID int `json:"userId"`
	DeviceID string `json:"deviceId"`
}

// Logout handler for POST /auth/logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request", http.StatusMethodNotAllowed) // Only allow POST
		return
	}

	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest) // Error decoding JSON body
		return
	}

	if err := db.RevokeDevice(req.DeviceID); err != nil {
		http.Error(w, "Could not revoke session", http.StatusInternalServerError) // Error revoking device
		return
	}

	// Clear the access token cookie
	http.SetCookie(w, &http.Cookie{
		Name: "access_token",
		Value: "",
		HttpOnly: true,
		Secure: !isDevelopment,
		SameSite: http.SameSiteStrictMode,
		Path: "/",
		MaxAge: -1,
	})

	// Best effort logging, log even in logout failures
	_ = db.LogActivity("logout", "", strconv.Itoa(req.UserID), req.DeviceID, r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out"))
}

// Pin request structure
type verifyPinRequest struct {
	PIN string `json:"pin"`
	DeviceID string `json:"deviceId"`
}

// VerifyPinHandler for POST /auth/verify-pin. X-User-ID in header, so token does not need to be reparsed to find userID
func VerifyPinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request", http.StatusMethodNotAllowed)
		return
	}

	var req verifyPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// Read X-User-ID
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Hashed pin
	hashedPin, err := db.GetPinByUserID(userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if hashedPin == "" {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare pin
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPin), []byte(req.PIN)); err != nil {
		// Wrong pin
		_ = db.LogActivity("pin_failed", "", strconv.Itoa(userID), req.DeviceID, r.RemoteAddr)
		http.Error(w, "Invalid PIN", http.StatusUnauthorized)
		return
	}

	// Finally success
	_ = db.LogActivity("pin_verified", "", strconv.Itoa(userID), req.DeviceID, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("PIN verified"))
}

// Me handler, GET auth/me for cookie check
func MeHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Header.Get("X-User-ID")

    w.Header().Set("Content-Type", "application/json")

    json.NewEncoder(w).Encode(struct {
        UserID string `json:"userId"`
    }{
        UserID: userID,
    })
}