// Package server contains the server part of Evento.
package server

import (
	"cmp"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	conn        *pgxpool.Pool
	databaseURL = cmp.Or(os.Getenv("DATABASE_URL"), "postgres://postgres@localhost:5432/evento")
)

func Build() (*http.ServeMux, error) {
	pconfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}

	conn, err = pgxpool.NewWithConfig(context.Background(), pconfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}

	// Start the server
	server := http.NewServeMux()
	server.HandleFunc("GET /available", available)
	server.HandleFunc("POST /reserve/naive", naive)
	server.HandleFunc("POST /reserve/safe", safe)

	return server, nil
}
