package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
)

func New(databaseConfig config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(databaseConfig.Driver, databaseConfig.URL)
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}

	db.SetMaxOpenConns(databaseConfig.MaxOpenConns)
	db.SetMaxIdleConns(databaseConfig.MaxIdleConns)
	db.SetConnMaxLifetime(databaseConfig.ConnMaxLifetime)
	db.SetConnMaxIdleTime(databaseConfig.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), databaseConfig.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return db, nil
}
