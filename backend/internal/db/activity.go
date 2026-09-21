// Log db activity

package db

import "time"

// Activity struct to display logs
type Activity struct {
    ID        int       `json:"id"`
    EventType string    `json:"eventType"`
    FilePath  string    `json:"filePath"`
    UserID    string    `json:"userId"`
    DeviceID  string    `json:"deviceId"`
    IPAddress string    `json:"ipAddress"`
    CreatedAt time.Time `json:"createdAt"`
}

// LogActivity inserts a new row into the activity_log table from auth.go
func LogActivity(eventType, filePath, userID, deviceID, ipAddress string) error {
	_, err := DB.Exec(
		`INSERT INTO activity_log (event_type, file_path, user_id, device_id, ip_address, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		eventType, filePath, userID, deviceID, ipAddress, time.Now(),
	)
	return err
}

func GetActivityLogs(limit int) ([]Activity, error) {
	rows, err := DB.Query(
        `SELECT id, event_type, file_path, user_id, device_id, ip_address, created_at
         FROM activity_log
         ORDER BY created_at DESC
         LIMIT ?`,
        limit,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    activities := []Activity{}
    for rows.Next() {
        var a Activity
        if err := rows.Scan(&a.ID, &a.EventType, &a.FilePath, &a.UserID, &a.DeviceID, &a.IPAddress, &a.CreatedAt); err != nil {
            return nil, err
        }
        activities = append(activities, a)
    }
    return activities, rows.Err()
}