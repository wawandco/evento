package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Show available rooms for an event
func available(w http.ResponseWriter, r *http.Request) {
	fmt.Println("info: available called")

	eventID := r.URL.Query().Get("event_id")
	query := `
	SELECT hotel_id, (contracted-(reserved+locked)) as available
	FROM event_hotel_rooms
	WHERE
		(reserved + locked) < contracted
		AND event_id = $1
	`

	rows, err := conn.Query(r.Context(), query, eventID)
	if err != nil {
		fmt.Println("error querying database", err)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type hotelAvailability struct {
		HotelID   string `json:"hotel_id"`
		Available int    `json:"available"`
	}

	var availability []hotelAvailability
	for rows.Next() {
		var hotel hotelAvailability

		if err := rows.Scan(&hotel.HotelID, &hotel.Available); err != nil {
			fmt.Println("error scanning row")
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}

		availability = append(availability, hotel)
	}

	dat, err := json.Marshal(availability)
	if err != nil {
		fmt.Println("error marshalling json", err)
		http.Error(w, "Error marshalling json", http.StatusInternalServerError)
		return
	}

	w.Write(dat)
}
