package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// schema for the database, it starts dropping existing tables to start fresh
var schema = `
	DROP TABLE IF EXISTS reservations;
	DROP TABLE IF EXISTS event_hotel_rooms;
	DROP TABLE IF EXISTS hotel;
	DROP TABLE IF EXISTS events;

	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS  events (
	    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	    name TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS hotel (
	    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	    name TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS event_hotel_rooms  (
	    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	    event_id UUID NOT NULL,
	    hotel_id UUID NOT NULL,

	    contracted INTEGER NOT NULL DEFAULT 0,
	    locked INTEGER NOT NULL DEFAULT 0,
	    reserved INTEGER NOT NULL DEFAULT 0,
	    updated_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS reservations (
	    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	    email TEXT NOT NULL,
	    event_hotel_rooms_id UUID,
	    number_of_rooms INTEGER
	);
`

// setup schema runs the schema SQL against the database using the
// provided connection pool.
func setupSchema(con *pgxpool.Pool) error {
	_, err := con.Exec(context.Background(), schema)
	if err != nil {
		return fmt.Errorf("error running schema: %w", err)
	}

	return nil
}
