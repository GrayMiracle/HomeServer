// Config

// config package for config files
package config

// imports
import (
	"encoding/json" // json encoding for config
	"os" // os file handling
)

// config struct similar to config.json
type Config struct {
	Port string `json:"port"`
	DefaultFolder string `json:"default_folder"`
	Domain string `json:"domain"`
	IsDevelopment bool `json:"IsDevelopment"`
	NoIPUsername string `json:"noip_username"`
	NoIPPassword string `json:"noip_password"`
	FrontendURL string `json:"frontend_url"`
}

// Load func to get config.json with the Config structu
func Load(path string) (*Config, error) {
	// os.ReadFile to get whole file into memory as byte slice, works since config files are small enough. 
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// emprt config struct to unmarshal config.json bytes
	var cfg Config

	// Unmarshaling data to unpack raw config.JSON bytes to struct
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	// If running in docker, default_folder exists and is used as config
	if v := os.Getenv("DEFAULT_FOLDER"); v != "" {
		cfg.DefaultFolder = v
	}

	// Return pointer to config file
	return &cfg, nil
}