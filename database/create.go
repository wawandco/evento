package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

// Creates the database to be used by Evento
func create(url string) error {
	// Connect to the default "postgres" database to create a new database
	defaultDB := strings.Replace(url, "evento", "postgres", 1)
	con, err := pgx.Connect(context.Background(), defaultDB)
	if err != nil {
		return fmt.Errorf("error connecting to the default database: %w", err)
	}

	defer con.Close(context.Background())

	// determine if the database already exists
	var exists bool
	err = con.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname='evento')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if database exists: %w", err)
	}

	if exists {
		fmt.Println("info: database 'evento' already exists")
		return nil
	}

	// Create the new database
	_, err = con.Exec(context.Background(), "CREATE DATABASE evento")
	if err != nil {
		return fmt.Errorf("error creating database: %w", err)
	}

	fmt.Println("info: database 'evento' created successfully")
	return nil
}
