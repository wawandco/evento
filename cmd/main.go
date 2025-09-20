package main

import (
	"cmp"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"evento/client"
	"evento/database"
	"evento/inventory"
	"evento/server"
)

// connection string to the database, defaults to a local Postgres instance
var databaseURL = cmp.Or(
	os.Getenv("DATABASE_URL"),
	"postgres://postgres@localhost:5432/evento",
)

func main() {
	modes := []string{"naive", "safe", "atomic"}
	usage := fmt.Sprintf("usage: evento [# of clients] [%s]", strings.Join(modes, "|"))

	args := os.Args
	if len(args) < 3 {
		fmt.Println(usage)
		return
	}

	clients, err := strconv.Atoi(args[1])
	if err != nil || clients <= 0 {
		fmt.Println("erro :invalid number of clients")
		fmt.Println(usage)
		return
	}

	mode := args[2]
	if !slices.Contains(modes, mode) {
		fmt.Println(usage)
		return
	}

	conn, err := database.Connect(databaseURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// Create the database and the schema
	err = database.Setup(conn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	fmt.Println("Database Ready")
	fmt.Println("Starting inventory")
	inventory.Print(conn)
	fmt.Println()

	// Run the server
	go func() {
		srv, err := server.New(conn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		fmt.Println("Server running")
		err = http.ListenAndServe(":8080", srv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
	}()

	wg := sync.WaitGroup{}
	for i := range clients {
		wg.Go(func() {
			client.Run(
				mode,
				fmt.Sprintf("client-%d", i),
				"7f3535c6-d5cb-44f0-b89b-4b349f01e49d",
			)
		})
	}

	fmt.Println("All clients started")

	// Wait for all clients to finish
	wg.Wait()
	fmt.Printf("\nFinal Inventory\n")

	// Print the inventory
	inventory.Print(conn)
}
