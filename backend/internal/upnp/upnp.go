// upnp port handling

package upnp

import (
	"log" // Logging
	"net" // Network operations
	"github.com/huin/goupnp/dcps/internetgateway2" // UPnP WANIP connection
)

var client *internetgateway2.WANIPConnection1

// Port open fn
func Open() error {
	// Discover UPnP devices
	clients, _, err := internetgateway2.NewWANIPConnection1Clients()
	// If theres an error or no clients return error
	if err != nil || len(clients) == 0 {
		return err
	}

	// Set first client found
	client = clients[0]

	// Get ip
	ip, err := getLocalIP()
	if err != nil {
		return err
	}

	// 
	for _, port := range []uint16{80, 443} {
		err := client.AddPortMapping("", port, "TCP", port, ip, true, "HomeServer", 0)
		if err != nil {
			log.Printf("UPnP failed to open port %d: %v", port, err)
		} else {
			log.Printf("UPnP opened port at %d -> %s", port, ip)
		}
	}
	return nil
}

// Close port
func Close() {
    if client == nil {
        return
    }
    for _, port := range []uint16{80, 443} {
        err := client.DeletePortMapping("", port, "TCP")
        if err != nil {
            log.Printf("Failed to close port! %d: %v", port, err)
        } else {
            log.Printf("UPnP closed port successfully: %d", port)
        }
    }
}

// Get local IP address
func getLocalIP() (string, error) {
	// Dial UDP to 8.8.8.8:80 to get local IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	// Close connection and return local IP
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}