package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
)

func New(databaseConfig config.DatabaseConfig) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseConfig.URL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	config.MaxConns = int32(databaseConfig.MaxOpenConns)
	config.MinConns = int32(databaseConfig.MaxIdleConns)
	config.MaxConnLifetime = databaseConfig.ConnMaxLifetime
	config.MaxConnIdleTime = databaseConfig.ConnMaxIdleTime
	config.ConnConfig.ConnectTimeout = databaseConfig.ConnectTimeout

	ctx, cancel := context.WithTimeout(context.Background(), databaseConfig.ConnectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return pool, nil
}
