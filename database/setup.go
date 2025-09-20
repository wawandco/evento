// Package database takes care of database operations like
package database

import (
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(conn *pgxpool.Pool, rooms int) error {
	err := create(conn.Config().ConnString())
	if err != nil {
		return err
	}

	err = createSchema(conn)
	if err != nil {
		return err
	}

	err = seed(conn, rooms)
	if err != nil {
		return err
	}

	return nil
}
