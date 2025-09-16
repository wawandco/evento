// Package server contains the server part of Evento.
package server

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool

// New server instance
func New(cx *pgxpool.Pool) (*http.ServeMux, error) {
	// Store the connection pool in a package variable
	conn = cx

	// Start the server
	server := http.NewServeMux()
	server.HandleFunc("GET /available", available)
	server.HandleFunc("POST /reserve/naive", naive)
	server.HandleFunc("POST /reserve/safe", safe)

	return server, nil
}
