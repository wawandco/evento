package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

// naive method for reserving rooms for an event. We use the database connection
// pool to query the database, update availability and create the reservation.
func naive(w http.ResponseWriter, r *http.Request) {
	// parse the event_id, hotel_id, email and number of rooms from the URL path and query parameters
	// this is done for simplicity, in a real application you would use a JSON body or form values
	eventID := r.PathValue("event_id")
	hotelID := r.PathValue("hotel_id")
	email := r.URL.Query().Get("email")

	// parse the number of rooms to reserve and validate is a positive integer
	rooms, err := strconv.Atoi(r.URL.Query().Get("rooms"))
	if err != nil || rooms <= 0 {
		w.Write([]byte("invalida number of rooms"))

		http.Error(w, "Invalid number of rooms", http.StatusBadRequest)
		return
	}

	// check if the number of rooms specified is available
	// in the event_hotel_rooms table availability
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
	err = conn.QueryRow(r.Context(), query, eventID, hotelID, rooms).Scan(&available)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Error querying availability", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Not enough rooms available", http.StatusConflict)
		return
	}

	// Update the availability by increasing the reserved rooms
	query = `
		UPDATE event_hotel_rooms
		SET reserved = reserved + $1
		WHERE
			event_id = $2
		AND
			hotel_id = $3
	`

	_, err = conn.Exec(r.Context(), query, rooms, eventID, hotelID)
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

	_, err = conn.Exec(r.Context(), query, eventID, hotelID, email, rooms)
	if err != nil {
		http.Error(w, "Error creating reservation", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("reservation successful"))
}
