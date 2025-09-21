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
	server.HandleFunc("GET /{event_id}/available", available)
	server.HandleFunc("POST /{event_id}/{hotel_id}/reserve/naive", naive)
	server.HandleFunc("POST /{event_id}/{hotel_id}/reserve/pessimistic", pessimistic)
	server.HandleFunc("POST /{event_id}/{hotel_id}/reserve/atomic", atomic)
	server.HandleFunc("POST /{event_id}/{hotel_id}/reserve/optimistic", optimistic)

	return server, nil
}
