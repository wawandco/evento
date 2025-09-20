package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Seed the database with initial data for testing purposes, it starts by cleaning up any existing data.
// rooms is the number of rooms to seed the database with per hotel.
func seed(con *pgxpool.Pool, rooms int) error {
	_, err := con.Exec(context.Background(), `
		DELETE FROM events;
		DELETE FROM hotel;
		DELETE FROM event_hotel_rooms;
		DELETE FROM reservations;
	`)
	if err != nil {
		return fmt.Errorf("error running data script: %w", err)
	}

	_, err = con.Exec(context.Background(), `
		INSERT INTO
		    events (id, name)
		VALUES
		    ($1, 'Evento')
		ON CONFLICT (id) DO NOTHING;`,
		"7f3535c6-d5cb-44f0-b89b-4b349f01e49d",
	)
	if err != nil {
		return fmt.Errorf("error running data script: %w", err)
	}

	// Sample data for hotels and event_hotel_rooms
	_, err = con.Exec(context.Background(), `
		INSERT INTO hotel (id, name)
		VALUES
		    ('1c3f3e2a-3f4e-4b62-8f4e-2b3e4f6c7d8e', 'Hotel A'),
		    ('2d4e5f6a-4b5e-6c7d-8e9f-0a1b2c3d4e5f', 'Hotel B'),
		    ('3e5f6a7b-5c6d-7e8f-9a0b-1c2d3e45f6aa', 'Hotel C')
		ON CONFLICT (id) DO NOTHING;`,
	)
	if err != nil {
		return fmt.Errorf("error inserting hotels: %w", err)
	}

	_, err = con.Exec(context.Background(), `
		INSERT INTO
		    event_hotel_rooms (id, event_id, hotel_id, contracted, locked, reserved)
		VALUES
		    ('a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', $1, '1c3f3e2a-3f4e-4b62-8f4e-2b3e4f6c7d8e', $2, 0, 0),
		    ('b2c3d4e5-f6a7-8b9c-0d1e-2f34a5b6c7da', $1, '2d4e5f6a-4b5e-6c7d-8e9f-0a1b2c3d4e5f', $2, 0, 0),
		    ('c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', $1, '3e5f6a7b-5c6d-7e8f-9a0b-1c2d3e45f6aa', $2, 0, 0)
		ON CONFLICT (id) DO NOTHING;`,
		"7f3535c6-d5cb-44f0-b89b-4b349f01e49d",
		rooms,
	)
	if err != nil {
		return fmt.Errorf("error inserting event hotel rooms: %w", err)
	}

	return nil
}
