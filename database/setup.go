// Package database takes care of database operations like
package database

import (
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Setup the database, creating it if it doesn't exist, creating the schema and seeding it with data.
// rooms is the number of rooms to seed the database with per hotel.
func Setup(conn *pgxpool.Pool, rooms int) error {
	err := create(conn.Config().ConnString())
	if err != nil {
		return err
	}

	err = setupSchema(conn)
	if err != nil {
		return err
	}

	err = seed(conn, rooms)
	if err != nil {
		return err
	}

	return nil
}
