package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
)

//go:embed schema.sql
var schema string

func createSchema() error {
	con, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %w", err)
	}

	defer con.Close(context.Background())

	_, err = con.Exec(context.Background(), schema)
	if err != nil {
		return fmt.Errorf("error running schema: %w", err)
	}

	fmt.Println("info: schema completed")

	return nil
}
