package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Databases map[string]DatabaseConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Kind     string
	URI      string
	Database string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: port,
		},
		Databases: map[string]DatabaseConfig{
			"primary": {
				Kind:     getEnv("DB_PRIMARY_KIND", "mongodb"),
				URI:      getEnv("DB_PRIMARY_URI", "mongodb://localhost:27017"),
				Database: getEnv("DB_PRIMARY_DATABASE", "appdb"),
			},
		},
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
