// Package client will take care of the client side of Evento.
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Run is a function that will call the server to reserve seats.
// It checks availability and based on it will call the reserve endpoint.
// It logs the result of the reservation. It will run until there
// is no availability in any hotel.
func Run(port, kind, ID, eventID string) {
	for {
		// Make a GET request to the Availability endpoint which will return a
		// JSON array of available hotels with available rooms.
		// Example: [{"hotel_id":"1","available":10},{"hotel_id":"2","available":5}]
		res, err := http.Get(fmt.Sprintf("http://localhost:%s/%s/available", port, eventID))
		if err != nil {
			continue
		}

		// Retries on non-200 status codes
		if res.StatusCode != http.StatusOK {
			continue
		}

		// Read the body of the response and retry.
		bb, err := io.ReadAll(res.Body)
		if err != nil {
			continue
		}

		// available is a struct to unmarshal the JSON response from the server
		// when checking availability.
		// Example: [{"hotel_id":"1","available":10},{"hotel_id":"2","available":5}]
		availability := []struct {
			HotelID   string `json:"hotel_id"`
			Available int    `json:"available"`
		}{}

		// Unmarshal the JSON response and retry.
		err = json.Unmarshal(bb, &availability)
		if err != nil {
			fmt.Println("client: error unmarshalling")
			continue
		}
		res.Body.Close()

		// If there is no availability, stop the client.
		if len(availability) == 0 {
			break
		}

		// random number of rooms 1-5
		rooms := int(time.Now().UnixNano()%5) + 1

		// If there is availability reserve 1 room in the first hotel with availability.
		url := fmt.Sprintf(
			"http://localhost:%s/%s/%s/reserve/%s?rooms=%d&email=%s@client.com",
			port, eventID, availability[0].HotelID, kind, rooms, ID,
		)

		_, err = http.Post(url, "application/x-www-form-urlencoded", nil)
		if err != nil {
			continue
		}

		// Indpendent of the response take a random sleep of less than 300 ms after calling reservation.
		// Using the mod operator on the current time in nanoseconds to generate a pseudo-random number.
		rand := time.Now().UnixNano() % 300
		time.Sleep(time.Duration(rand) * time.Millisecond)
	}
}
