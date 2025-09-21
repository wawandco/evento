package server

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

// Optimistic method for reserving rooms for an event, it uses a transaction
// with optimistic locking based on updated_at timestamp to prevent overbooking.
func optimistic(w http.ResponseWriter, r *http.Request) {
	// parse the event_id, hotel_id, email and number of rooms from the URL path and query parameters
	eventID := r.PathValue("event_id")
	hotelID := r.PathValue("hotel_id")
	email := r.URL.Query().Get("email")

	// parse the number of rooms to reserve and validate is a positive integer
	rooms, err := strconv.Atoi(r.URL.Query().Get("rooms"))
	if err != nil || rooms <= 0 {
		w.Write([]byte("invalid number of rooms"))

		http.Error(w, "Invalid number of rooms", http.StatusBadRequest)
		return
	}

	// Start a new transaction with the request context
	tx, err := conn.Begin(r.Context())
	if err != nil {
		http.Error(w, "error starting transaction", http.StatusInternalServerError)
		return
	}
	// Defer a rollback in case of any errors.
	defer tx.Rollback(r.Context())

	// check if quantity is available and get current updated_at timestamp
	query := `
		SELECT updated_at
		FROM event_hotel_rooms
		WHERE
			event_id = $1
		AND
			hotel_id = $2
		AND
			contracted - (reserved + locked) >= $3
	`

	var updatedAt time.Time
	err = tx.QueryRow(r.Context(), query, eventID, hotelID, rooms).Scan(&updatedAt)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Error querying availability:"+err.Error(), http.StatusInternalServerError)
		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Not enough rooms available", http.StatusConflict)
		return
	}

	// update rooms availability by increasing the reserved rooms and updating timestamp
	query = `
		UPDATE event_hotel_rooms
		SET reserved = reserved + $1, updated_at = NOW()
		WHERE
			event_id = $2
		AND
			hotel_id = $3
		AND
			updated_at = $4
	`

	res, err := tx.Exec(r.Context(), query, rooms, eventID, hotelID, updatedAt)
	if err != nil {
		http.Error(w, "Error reserving rooms", http.StatusInternalServerError)
		return
	}

	// Check if the update affected any rows (optimistic locking)
	if res.RowsAffected() == 0 {
		// Log the conflict for debugging
		http.Error(w, "Conflict: data was modified by another transaction", http.StatusConflict)
		return
	}

	// insert the reservation record in the reservations table
	query = `
		INSERT INTO
			reservations (event_hotel_rooms_id, email, number_of_rooms)
		VALUES (
			(SELECT id FROM event_hotel_rooms WHERE event_id = $1 AND hotel_id = $2),
			$3,
			$4
		);
	`

	_, err = tx.Exec(r.Context(), query, eventID, hotelID, email, rooms)
	if err != nil {
		http.Error(w, "Error creating reservation", http.StatusInternalServerError)
		return
	}

	tx.Commit(r.Context())
	w.Write([]byte("reservation successful"))
}
