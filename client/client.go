// Package client will take care of the client side of Evento.
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// A client is a goroutine that will call the server to reserve seats.
// It will check availability and based on it will call the reserve endpoint.
// It will log the result of the reservation.

func Client(eventID string) {
	type available struct {
		HotelID   string `json:"hotel_id"`
		Available int    `json:"available"`
	}

	for {
		availability := []available{}
		// Make a GET request to the Availability endpoint
		res, err := http.Get("http://localhost:8080/available")
		if err != nil {
			fmt.Println("client: error requesting")
			continue
		}

		if res.StatusCode != http.StatusOK {
			fmt.Println("client: error code")
			continue
		}

		bb, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("client: error reading body")
			continue
		}

		err = json.Unmarshal(bb, &availability)
		if err != nil {
			fmt.Println("client: error unmarshalling")
			continue
		}

		res.Body.Close()
		if len(availability) == 0 {
			fmt.Println("No availability")
			break
		}

		url := fmt.Sprintf("http://localhost:8080/reserve?event_id=%s&hotel_id=%s&rooms=1", eventID, availability[0].HotelID)
		resp, err := http.Post(url, "application/x-www-form-urlencoded", nil)
		if err != nil {
			fmt.Println("client: error reserving")
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println("client: error reserving")
			continue
		}

		fmt.Println("client: reservation successful")
	}
}
