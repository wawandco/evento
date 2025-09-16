package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"evento/client"
	"evento/database"
	"evento/results"
	"evento/server"
)

func main() {
	usage := "usage: evento <naive|safe> clients"

	args := os.Args
	if len(args) < 3 {
		fmt.Println(usage)
		return
	}

	clients, err := strconv.Atoi(args[1])
	if err != nil || clients <= 0 {
		fmt.Println(usage)
		return
	}

	mode := args[2]
	if mode != "naive" && mode != "safe" {
		fmt.Println(usage)
		return
	}

	// Create the database and the schema
	err = database.Setup()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// Load initial data
	err = database.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// Run the server
	go func() {
		srv, err := server.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		fmt.Println("info: server running on :8080")
		err = http.ListenAndServe(":8080", srv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
	}()

	fmt.Printf("Running %d concurrent clients until inventory has been depleted.", clients)
	wg := sync.WaitGroup{}
	for i := range clients {
		wg.Go(func() {
			fmt.Printf("client %d\n", i)
			client.Run(
				mode,
				fmt.Sprintf("client-%d", i),
				"7f3535c6-d5cb-44f0-b89b-4b349f01e49d",
			)
		})
	}

	// Wait for all clients to finish
	wg.Wait()
	fmt.Println("Rooms Inventory reserved")

	// Print the inventory
	results.Print()
}
