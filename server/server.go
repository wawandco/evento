// Package server contains the server part of Evento.
package server

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool

func Build() (*http.ServeMux, error) {
	pconfig, err := pgxpool.ParseConfig("postgres://postgres@localhost:5432/evento")
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
	server.HandleFunc("POST /reserve", reserve)

	return server, nil
}
