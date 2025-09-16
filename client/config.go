package client

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerAddress string `json:"server_address"`
	ServerPort    int    `json:"server_port"`
	UseTLS        bool   `json:"use_tls"`
	ServerName    string `json:"server_name"`
	CertPath      string `json:"cert_path"`
}

func LoadConfig() (*Config, error) {
	// Use project-relative paths instead of home directory
	configPath := "config.json"

	// Create default config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := &Config{
			ServerAddress: "your-server",
			ServerPort:    8080,
			UseTLS:        true,
			ServerName:    "server",
			CertPath:      "certs/server.crt",
		}

		if err := SaveConfig(defaultConfig); err != nil {
			return nil, err
		}

		return defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	// Save config in project root
	configPath := "config.json"
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
