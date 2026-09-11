package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment string
	Port        string
	Database    DatabaseConfig
	Logger      LoggerConfig
}

type LoggerConfig struct {
	Level  string
	Format string
}

type DatabaseConfig struct {
	Driver          string
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	ConnectTimeout  time.Duration
}

func Load() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	databaseConfig, err := loadDatabaseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	loggerConfig, err := loadLoggerConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Environment: valueOrDefault("APP_ENV", "development"),
		Port:        valueOrDefault("PORT", "8080"),
		Database:    databaseConfig,
		Logger:      loggerConfig,
	}, nil
}

func loadDatabaseConfig(databaseURL string) (DatabaseConfig, error) {
	maxOpenConns, err := intValue("DATABASE_MAX_OPEN_CONNS", 10, 1)
	if err != nil {
		return DatabaseConfig{}, err
	}

	maxIdleConns, err := intValue("DATABASE_MAX_IDLE_CONNS", 5, 0)
	if err != nil {
		return DatabaseConfig{}, err
	}
	if maxIdleConns > maxOpenConns {
		return DatabaseConfig{}, fmt.Errorf("DATABASE_MAX_IDLE_CONNS cannot exceed DATABASE_MAX_OPEN_CONNS")
	}

	maxLifetime, err := durationValue("DATABASE_CONN_MAX_LIFETIME", 30*time.Minute)
	if err != nil {
		return DatabaseConfig{}, err
	}

	maxIdleTime, err := durationValue("DATABASE_CONN_MAX_IDLE_TIME", 5*time.Minute)
	if err != nil {
		return DatabaseConfig{}, err
	}

	connectTimeout, err := durationValue("DATABASE_CONNECT_TIMEOUT", 5*time.Second)
	if err != nil {
		return DatabaseConfig{}, err
	}

	return DatabaseConfig{
		Driver:          valueOrDefault("DATABASE_DRIVER", "postgres"),
		URL:             databaseURL,
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: maxLifetime,
		ConnMaxIdleTime: maxIdleTime,
		ConnectTimeout:  connectTimeout,
	}, nil
}

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func intValue(key string, fallback, min int) (int, error) {
	value := valueOrDefault(key, strconv.Itoa(fallback))
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < min {
		return 0, fmt.Errorf("%s must be an integer greater than or equal to %d", key, min)
	}

	return parsed, nil
}

func durationValue(key string, fallback time.Duration) (time.Duration, error) {
	value := valueOrDefault(key, fallback.String())
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration, for example 30s or 5m", key)
	}

	return parsed, nil
}

var validLogLevels = map[string]bool{
	"debug": true, "info": true, "warn": true, "error": true,
}

var validLogFormats = map[string]bool{
	"json": true, "text": true,
}

func loadLoggerConfig() (LoggerConfig, error) {
	level := strings.ToLower(valueOrDefault("LOG_LEVEL", "info"))
	if !validLogLevels[level] {
		return LoggerConfig{}, fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}

	format := strings.ToLower(valueOrDefault("LOG_FORMAT", "json"))
	if !validLogFormats[format] {
		return LoggerConfig{}, fmt.Errorf("LOG_FORMAT must be one of: json, text")
	}

	return LoggerConfig{Level: level, Format: format}, nil
}
