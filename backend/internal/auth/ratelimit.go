// Rate limiter for pin

package auth

import (
	"net" // net for IP handling
	"net/http" // http for middleware
	"sync" // sync for mutex
	"time" // time for rate limit duration
)

// attemptRecord for failed login attempts per IP address
type attemptRecord struct {
	count int
	blockedAt time.Time
}

// Login limiter struct with mutex for concurrent access
var loginLimiter = struct {
	sync.Mutex
	attempts map[string]*attemptRecord
}{
	attempts: make(map[string]*attemptRecord),
}

// Attempts and block length
const (
	maxAttempts = 5
	blockDuration = 15 * time.Minute
)

// getIP to get the IP of the request
func getIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// Check whether IP is rate limited, true if blocked
func IsRateLimited(r *http.Request) bool {
	ip := getIP(r)

	loginLimiter.Lock()
	defer loginLimiter.Unlock()

	record, exists := loginLimiter.attempts[ip]
	// record does not exist
	if !exists {
		return false
	}

	// Otherwise check blockedAt time
	if !record.blockedAt.IsZero() {
		// If duration passed, reset ban
		if time.Since(record.blockedAt) > blockDuration {
			delete(loginLimiter.attempts, ip)
			return false
		}
		// If still banned
		return true
	}

	return false // not banned yet, attempts left
}

// Increment attempts
func RecordFailedAttempt(r *http.Request) {
	ip := getIP(r)

	// Lock for concurrent access
	loginLimiter.Lock()
	defer loginLimiter.Unlock()

	// Get or create record for IP
	record, exists := loginLimiter.attempts[ip]
	if !exists {
		record = &attemptRecord{}
		loginLimiter.attempts[ip] = record
	}

	// Increment count and check if should block
	record.count++
	if record.count >= maxAttempts {
		record.blockedAt = time.Now()
	}
}

// Successful login
func RecordSuccessfulLogin(r *http.Request) {
	ip := getIP(r)

	loginLimiter.Lock()
	defer loginLimiter.Unlock()

	delete(loginLimiter.attempts, ip)
}