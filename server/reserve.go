package server

import (
	"net/http"
	"strconv"
)

// Reserve rooms for an event
func reserve(w http.ResponseWriter, r *http.Request) {
	// Determine the event ID and HotelID from the request
	// Check availability in the database
	// If available, create a reservation record
	// Respond with success or failure

	eventID := r.URL.Query().Get("event_id")
	hotelID := r.URL.Query().Get("hotel_id")

	rooms, err := strconv.Atoi(r.URL.Query().Get("rooms"))
	if err != nil || rooms <= 0 {
		w.Write([]byte("invalida number of rooms"))

		http.Error(w, "Invalid number of rooms", http.StatusBadRequest)
		return
	}

	// check if quantity is available
	query := `
		SELECT (contracted - (reserved + locked)) as available
		FROM event_hotel_rooms
		WHERE event_id = $1 AND hotel_id = $2
	`

	var available int
	err = conn.QueryRow(r.Context(), query, eventID, hotelID).Scan(&available)
	if err != nil {
		http.Error(w, "Error querying availability", http.StatusInternalServerError)
		return
	}

	if available < rooms {
		http.Error(w, "Not enough rooms available", http.StatusConflict)
		return
	}

	// reserve the rooms
	query = `
		UPDATE event_hotel_rooms
		SET reserved = reserved + $1
		WHERE event_id = $2 AND hotel_id = $3
	`

	_, err = conn.Exec(r.Context(), query, rooms, eventID, hotelID)
	if err != nil {
		http.Error(w, "Error reserving rooms", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("reservation successful"))
}
