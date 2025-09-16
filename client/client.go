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
func Run(kind, ID, eventID string) {
	for {
		fmt.Printf("info: client %s requesting availability\n", ID)
		// Make a GET request to the Availability endpoint which will return a
		// JSON array of available hotels with available rooms.
		// Example: [{"hotel_id":"1","available":10},{"hotel_id":"2","available":5}]
		res, err := http.Get(fmt.Sprintf("http://localhost:8080/available?event_id=%s", eventID))
		if err != nil {
			fmt.Println("client: error requesting")
			continue
		}

		// Retries on non-200 status codes
		if res.StatusCode != http.StatusOK {
			fmt.Println("client: error code")
			continue
		}

		// Read the body of the response and retry.
		bb, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("client: error reading body")
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
			fmt.Printf("info: client %s found no availability, stopping \n", ID)
			break
		}

		// If there is availability reserve 1 room in the first hotel with availability.
		url := fmt.Sprintf(
			"http://localhost:8080/reserve/%s?event_id=%s&hotel_id=%s&rooms=1&email=%s@client.com",
			kind, eventID, availability[0].HotelID, ID,
		)

		resp, err := http.Post(url, "application/x-www-form-urlencoded", nil)
		if err != nil {
			fmt.Println("client: error reserving")
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println("client: error reserving")
			continue
		}

		fmt.Printf("info: client %s reservation successful\n", ID)

		// Random sleep of less than 300 ms after reservation. Using the current time
		// in nanoseconds to generate a pseudo-random number.
		rand := time.Now().UnixNano() % 300
		time.Sleep(time.Duration(rand) * time.Millisecond)
	}
}
