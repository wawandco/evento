package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

// Atomic method for reserving rooms for an event, it uses a transaction
// without locking. This can lead to overbooking to demonstrate the difference
// between thread-safe and atomic operations.
func atomic(w http.ResponseWriter, r *http.Request) {
	// parse the event_id, hotel_id, email and number of rooms from the URL path and query parameters
	// this is done for simplicity, in a real application you would use a JSON body or form values
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

	// Start a new transaction with the request context to avoid
	// "transaction has already been committed or rolled back" errors
	// if the client disconnects before the transaction is committed.
	tx, err := conn.Begin(r.Context())
	if err != nil {
		http.Error(w, "error starting transaction", http.StatusInternalServerError)
		return
	}
	// Defer a rollback in case of any errors.
	defer tx.Rollback(r.Context())

	// check if quantity is available WITHOUT a FOR UPDATE lock
	// This demonstrates that transactions alone don't prevent overbooking
	// because rows are not locked and other transactions can modify them
	// before the current transaction is committed.
	query := `
		SELECT true
		FROM event_hotel_rooms
		WHERE
			event_id = $1
		AND
			hotel_id = $2
		AND
			contracted - (reserved + locked) >= $3
	`

	var available bool
	err = tx.QueryRow(r.Context(), query, eventID, hotelID, rooms).Scan(&available)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Error querying availability", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Not enough rooms available", http.StatusConflict)
		return
	}

	// update rooms availability by increasing the reserved rooms in the
	// event_hotel_rooms table
	query = `
		UPDATE event_hotel_rooms
		SET reserved = reserved + $1
		WHERE
			event_id = $2
		AND
			hotel_id = $3
	`

	_, err = tx.Exec(r.Context(), query, rooms, eventID, hotelID)
	if err != nil {
		http.Error(w, "Error reserving rooms", http.StatusInternalServerError)
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
