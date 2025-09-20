package main

import (
	"cmp"
	"flag"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"evento/client"
	"evento/database"
	"evento/inventory"
	"evento/server"
)

var (
	mode    string
	clients int
	rooms   int
	servers int
)

func init() {
	flag.StringVar(&mode, "mode", "naive", "mode of operation: naive, safe, atomic, optimistic, defaults to naive")
	flag.IntVar(&clients, "clients", 200, "number of concurrent clients, defaults to 200")
	flag.IntVar(&rooms, "rooms", 200, "number of rooms per hotel, defaults to 200")
	flag.IntVar(&servers, "servers", 2, "number of server instances, defaults to 2")
	flag.Parse()
}

func main() {
	modes := []string{"naive", "atomic", "pessimistic", "optimistic"}
	usage := fmt.Sprintf("usage: evento [# of clients] [%s] [-rooms number]", strings.Join(modes, "|"))

	if clients <= 0 {
		fmt.Println("error: invalid number of clients")
		fmt.Println(usage)
		return
	}

	if !slices.Contains(modes, mode) {
		fmt.Println(usage)
		return
	}

	dbURL := cmp.Or(os.Getenv("DATABASE_URL"), "postgres://postgres@localhost:5432/evento")
	conn, err := database.Connect(dbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// Create the database and the schema
	err = database.Setup(conn, rooms)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	fmt.Println("- Database Ready (created, migrated and seeded)")
	fmt.Printf("- Using %d rooms per hotel\n", rooms)
	fmt.Println("> Starting inventory")
	inventory.Print(conn)
	fmt.Println()

	// Run the servers
	runningPorts := []string{}
	for v := range servers {
		port := 8080 + v
		go func(port int) {
			srv, err := server.New(conn)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				return
			}

			fmt.Printf("- Server running on port %d\n", port)

			// Start the server
			err = http.ListenAndServe(fmt.Sprintf(":%d", port), srv)
			if err != nil {
				return
			}
		}(port)
		runningPorts = append(runningPorts, fmt.Sprintf("%d", port))
	}

	// Start timing the execution
	start := time.Now()

	wg := sync.WaitGroup{}
	for i := range clients {
		wg.Go(func() {
			client.Run(
				// Randomly pick a server port to distribute the load
				runningPorts[time.Now().UnixNano()%int64(len(runningPorts))],
				mode,
				fmt.Sprintf("client-%d", i),
				"7f3535c6-d5cb-44f0-b89b-4b349f01e49d",
			)
		})
	}

	fmt.Printf("- %d `%s` clients making reservations\n", clients, mode)

	// Wait for all clients to finish
	wg.Wait()

	// Calculate and display execution time
	duration := time.Since(start)
	fmt.Printf("\n> Execution time: %v\n", duration)

	fmt.Printf("\n> Final Inventory\n")

	// Print the inventory
	inventory.Print(conn)
}
