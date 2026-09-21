// All devices user uses

// DB package for db files
package db

import (
	"time" // time for device expiration
	"golang.org/x/crypto/bcrypt" // bcrypt for hashing refresh tokens
	"database/sql" // sql database handling
)

// One row in the devices table — one "logged in session" holding a refresh token.
type Device struct {
	ID               int
	DeviceID         string
	DeviceName       string
	UserID           int
	RefreshTokenHash string
	Trusted          bool
	ExpiresAt        time.Time
}

// CreateDevice inserts a new device row, or updates existing row if already logged in before with an existing device_id
func CreateDevice(deviceID, deviceName string, userID int, refreshTokenHash string, expiresAt time.Time) error {
	_, err := DB.Exec(
		`INSERT INTO devices (device_id, device_name, user_id, refresh_token_hash, expires_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(device_id) DO UPDATE SET
			refresh_token_hash = excluded.refresh_token_hash,
			expires_at = excluded.expires_at,
			user_id = excluded.user_id`,
		deviceID, deviceName, userID, refreshTokenHash, expiresAt,
	)
	return err
}

// GetDevicesByUser fetches every device row for a user.
func GetDevicesByUser(userID int) ([]Device, error) {
	rows, err := DB.Query(
		`SELECT id, device_id, device_name, user_id, refresh_token_hash, trusted, expires_at
		 FROM devices WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.ID, &d.DeviceID, &d.DeviceName, &d.UserID,
			&d.RefreshTokenHash, &d.Trusted, &d.ExpiresAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

// GetDeviceByID fetches exactly one device row by its device_id.
func GetDeviceByID(deviceID string) (*Device, error) {
	var d Device
	row := DB.QueryRow(
		`SELECT id, device_id, device_name, user_id, refresh_token_hash, trusted, expires_at
		 FROM devices WHERE device_id = ?`,
		deviceID,
	)
	err := row.Scan(&d.ID, &d.DeviceID, &d.DeviceName, &d.UserID,
		&d.RefreshTokenHash, &d.Trusted, &d.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, nil // not found, not a real error
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// VerifyRefreshToken checks a raw token against stored hash for the singular device row. Identify with userId -> deviceId. Returns device if matches, nil if not (invalid/expired/doesnt exist)
func VerifyRefreshToken(userID int, deviceID, rawToken string) (*Device, error) {
	d, err := GetDeviceByID(deviceID)
	if err != nil {
		return nil, err
	}

	// Device must exist and belong to user
	if d == nil || d.UserID != userID {
		return nil, nil 
	}

	// Hash must match
	if err := bcrypt.CompareHashAndPassword([]byte(d.RefreshTokenHash), []byte(rawToken)); err != nil {
		return nil, nil
	}

	// Must not be expired
	if time.Now().After(d.ExpiresAt) {
		return nil, nil
	}

	return d, nil
}

// RevokeDevice deletes a device row — used for logout.
func RevokeDevice(deviceID string) error {
	_, err := DB.Exec("DELETE FROM devices WHERE device_id = ?", deviceID)
	return err
}

// UpdateLastSeen bumps last_seen — call this on every successful refresh.
func UpdateLastSeen(deviceID string) error {
	_, err := DB.Exec(
		"UPDATE devices SET last_seen = ? WHERE device_id = ?",
		time.Now(), deviceID,
	)
	return err
}