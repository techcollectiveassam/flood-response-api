package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
)

func New(databaseConfig config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(databaseConfig.Driver, databaseConfig.URL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
