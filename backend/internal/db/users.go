// Users

// DB package for db files
package db

// imports
import (
	"database/sql" // sql database handling
)

// Detect if user exists
func UserExists() (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create user with user and pass
func CreateUser(username, hashedPassword, hashedPin string) error {
	_, err := DB.Exec(
		"INSERT INTO users (username, password, pin) VALUES (?, ?, ?)",
		username, hashedPassword, hashedPin,
	)
	return err
}

// Get user by username for login
func GetUserByUsername(username string) (id int, hashedPassword string, err error) {
	row := DB.QueryRow(
		"SELECT id, password FROM users WHERE username = ?",
		username,
	)
	err = row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, "", nil // user not found, not an error
	}
	return id, hashedPassword, err
}

// Get pin by userid

func GetPinByUserID(userID int) (string, error) {
	var hashedPin string
	row := DB.QueryRow("SELECT pin FROM users WHERE id = ?", userID)
	err := row.Scan(&hashedPin)
	if err == sql.ErrNoRows {
		return "", nil // user not found
	}
	return hashedPin, err
}