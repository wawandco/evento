package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schema string

func createSchema(con *pgxpool.Pool) error {
	_, err := con.Exec(context.Background(), schema)
	if err != nil {
		return fmt.Errorf("error running schema: %w", err)
	}

	return nil
}
