// Package server contains the server part of Evento.
package server

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func Build() (*http.ServeMux, error) {
	cx, err := pgx.Connect(context.Background(), "postgres://postgres@localhost:5432/evento")
	if err != nil {
		return nil, err
	}

	conn = cx

	// Start the server
	server := http.NewServeMux()
	server.HandleFunc("GET /available", available)
	server.HandleFunc("POST /reserve", reserve)

	return server, nil
}
