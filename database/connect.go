package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect to the database and return a connection pool.
func Connect(url string) (*pgxpool.Pool, error) {
	pconfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DATABASE_URL: %w", err)
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), pconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	return conn, nil
}
