// Package database takes care of database operations like
package database

import (
	"cmp"
	_ "embed"
	"os"
)

// connection string to the database, defaults to a local Postgres instance
var databaseURL = cmp.Or(os.Getenv("DATABASE_URL"), "postgres://postgres@localhost:5432/evento")

func Setup() error {
	err := create()
	if err != nil {
		return err
	}

	err = createSchema()
	if err != nil {
		return err
	}

	return nil
}
