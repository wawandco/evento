package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed data.sql
var data string

func Load(con *pgxpool.Pool) error {
	_, err := con.Exec(context.Background(), data)
	if err != nil {
		return fmt.Errorf("error running data script: %w", err)
	}

	fmt.Println("info: data loaded")
	return nil
}
