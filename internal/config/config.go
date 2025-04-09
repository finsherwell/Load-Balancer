package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the load balancer configuration
type Config struct {
	Backends   []string `json:"backends"`
	Algorithm  string   `json:"algorithm"`
	MaxConns   int      `json:"max_connections"`
	HealthPath string   `json:"health_path"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}

	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	return &config, nil
}
