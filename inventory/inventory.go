// Package results allows to get the results of the reservations
// after running the clients against the server.
package inventory

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Print the results of the reservations after running the clients against the server.
func Print(conn *pgxpool.Pool) {
	query := `
		SELECT
			hotel.name,
			event_hotel_rooms.contracted,
			COALESCE(sum(reservations.number_of_rooms), 0) as reservations,
			event_hotel_rooms.reserved as rooms_reserved
		FROM
			hotel
			LEFT JOIN event_hotel_rooms ON hotel.ID = event_hotel_rooms.hotel_id
			LEFT JOIN reservations ON event_hotel_rooms.ID = reservations.event_hotel_rooms_id
		GROUP BY
			hotel.Name, event_hotel_rooms.reserved, event_hotel_rooms.contracted
		ORDER BY
			name, reservations DESC;
	`

	type hotelInventory struct {
		Name          string
		Contracted    string
		Reservations  int
		RoomsReserved int
	}

	data := []hotelInventory{}
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("error executing query: %s", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var row hotelInventory

		err := rows.Scan(&row.Name, &row.Contracted, &row.Reservations, &row.RoomsReserved)
		if err != nil {
			fmt.Printf("error scanning row: %s\n", err)
			return
		}

		data = append(data, row)
	}

	for _, v := range data {
		fmt.Printf("| Hotel: %s, Contracted: %s, Reservations: %d, Rooms Reserved: %d\n", v.Name, v.Contracted, v.Reservations, v.RoomsReserved)
	}
	fmt.Println("--------------------------------------")
}
