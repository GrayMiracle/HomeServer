// Dynamic DNS updates for noip

package noip

import (
	"fmt" // Printing
	"io" // I/O for HTTP response body
	"log" // Logging
	"net/http" // HTTP client
	"strings" // String manipulation
	"time" // Time for scheduling updates
)

const updateInterval = 5 * time.Minute // Update every 5 minutes

// Start func for updater goroutine
func Start(hostname, username, password string) {
	go func() {
        // Update immediately on startup
        if err := update(hostname, username, password); err != nil {
            log.Printf("No-IP update failed: %v", err)
        }
        // Then update every 5 minutes
        ticker := time.NewTicker(updateInterval)
        defer ticker.Stop()
        for range ticker.C {
            if err := update(hostname, username, password); err != nil {
                log.Printf("No-IP update failed: %v", err)
            }
        }
    }()
}

// update func to send current IP
func update(hostname, username, password string) error {
	url := fmt.Sprintf("https://dynupdate.no-ip.com/nic/update?hostname=%s", hostname)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

    // No-IP requires basic auth and a descriptive User-Agent
    req.SetBasicAuth(username, password)
    req.Header.Set("User-Agent", "HomeCloud/1.0 personal-homeserver")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    response := strings.TrimSpace(string(body))
    log.Printf("No-IP update response: %s", response)

    // No-IP returns "good <ip>" on success, "nochg <ip>" if IP hasn't changed
    if strings.HasPrefix(response, "good") || strings.HasPrefix(response, "nochg") {
        return nil
    }

    return fmt.Errorf("No-IP returned error: %s", response)
}