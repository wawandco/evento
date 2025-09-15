// Package results allows to get the results of the reservations
// after running the clients against the server.
package results

import (
	"cmp"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// connection string to the database, defaults to a local Postgres instance
var databaseURL = cmp.Or(os.Getenv("DATABASE_URL"), "postgres://postgres@localhost:5432/evento")

// Print the results of the reservations after running the clients against the server.
func Print() {
	con, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Printf("error connecting to the database: %s\n", err)
		return
	}

	defer con.Close(context.Background())

	query := `
	SELECT
		hotel.name,
		event_hotel_rooms.contracted,
		sum(reservations.number_of_rooms) as reservations,
		event_hotel_rooms.reserved as rooms_reserved
	FROM
		hotel
		LEFT JOIN event_hotel_rooms ON hotel.ID = event_hotel_rooms.hotel_id
		JOIN reservations ON event_hotel_rooms.ID = reservations.event_hotel_rooms_id
	GROUP BY
		hotel.Name, event_hotel_rooms.reserved, event_hotel_rooms.contracted
	ORDER BY
		name, reservations DESC;
	`

	type ResultingData struct {
		Name          string
		Contracted    string
		Reservations  int
		RoomsReserved int
	}

	data := []ResultingData{}
	rows, err := con.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("error executing query: %s", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var row ResultingData

		err := rows.Scan(&row.Name, &row.Contracted, &row.Reservations, &row.RoomsReserved)
		if err != nil {
			fmt.Printf("error scanning row: %s\n", err)
			return
		}

		data = append(data, row)
	}

	fmt.Printf("\nAfter Reservations\n")
	for _, v := range data {
		fmt.Printf("| Hotel: %s, Contracted: %s, Reservations: %d, Rooms Reserved: %d\n", v.Name, v.Contracted, v.Reservations, v.RoomsReserved)
	}
	fmt.Println("--------------------------------------")
}
