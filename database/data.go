package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
)

//go:embed data.sql
var data string

func Load() error {
	con, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %w", err)
	}

	defer con.Close(context.Background())

	_, err = con.Exec(context.Background(), data)
	if err != nil {
		return fmt.Errorf("error running data script: %w", err)
	}

	fmt.Println("info: data loaded")
	return nil
}
