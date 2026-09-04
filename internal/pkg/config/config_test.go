package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadUsesDatabaseDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@db:5432/flood_response?sslmode=disable")
	t.Setenv("DATABASE_DRIVER", "")
	t.Setenv("DATABASE_MAX_OPEN_CONNS", "")
	t.Setenv("DATABASE_MAX_IDLE_CONNS", "")
	t.Setenv("DATABASE_CONN_MAX_LIFETIME", "")
	t.Setenv("DATABASE_CONN_MAX_IDLE_TIME", "")
	t.Setenv("DATABASE_CONNECT_TIMEOUT", "")

	configuration, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if configuration.Database.Driver != "postgres" {
		t.Errorf("Driver = %q, want postgres", configuration.Database.Driver)
	}
	if configuration.Database.MaxOpenConns != 10 {
		t.Errorf("MaxOpenConns = %d, want 10", configuration.Database.MaxOpenConns)
	}
	if configuration.Database.MaxIdleConns != 5 {
		t.Errorf("MaxIdleConns = %d, want 5", configuration.Database.MaxIdleConns)
	}
	if configuration.Database.ConnMaxLifetime != 30*time.Minute {
		t.Errorf("ConnMaxLifetime = %s, want 30m", configuration.Database.ConnMaxLifetime)
	}
	if configuration.Database.ConnMaxIdleTime != 5*time.Minute {
		t.Errorf("ConnMaxIdleTime = %s, want 5m", configuration.Database.ConnMaxIdleTime)
	}
	if configuration.Database.ConnectTimeout != 5*time.Second {
		t.Errorf("ConnectTimeout = %s, want 5s", configuration.Database.ConnectTimeout)
	}
}

func TestLoadRejectsInvalidDatabasePoolSize(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@db:5432/flood_response?sslmode=disable")
	t.Setenv("DATABASE_MAX_OPEN_CONNS", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "DATABASE_MAX_OPEN_CONNS") {
		t.Errorf("Load() error = %q, want DATABASE_MAX_OPEN_CONNS", err)
	}
}
