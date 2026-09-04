package config

import (
	"fmt"
	"os"
)

type Config struct {
	Environment string
	Port        string
	Database    DatabaseConfig
}

type DatabaseConfig struct {
	Driver string
	URL    string
}

func Load() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return &Config{
		Environment: valueOrDefault("APP_ENV", "development"),
		Port:        valueOrDefault("PORT", "8080"),
		Database: DatabaseConfig{
			Driver: valueOrDefault("DATABASE_DRIVER", "postgres"),
			URL:    databaseURL,
		},
	}, nil
}

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
