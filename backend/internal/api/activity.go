// DB Log display route

package api

import (
	"encoding/json" // JSON encoding/decoding
	"net/http" // HTTP handling
	"homeserver/internal/db" // Database handling for logging
)

func LogHandler(w http.ResponseWriter, r *http.Request) {
	activities, err := db.GetActivityLogs(100)
	if err != nil {
		http.Error(w, "Couldn't fetch logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // JSON response header
	json.NewEncoder(w).Encode(activities)
}