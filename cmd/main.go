package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"evento/client"
	"evento/database"
	"evento/server"
)

func main() {
	// Create the database and the schema
	err := database.Setup()
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
	srv, err := server.Build()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	go func() {
		fmt.Println("info: server running on :8080")
		err = http.ListenAndServe(":8080", srv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
	}()

	time.Sleep(1 * time.Second)

	wg := sync.WaitGroup{}
	for i := range 200 {
		wg.Go(func() {
			fmt.Printf("client %d\n", i)
			client.Run(
				fmt.Sprintf("client-%d", i),
				"7f3535c6-d5cb-44f0-b89b-4b349f01e49d",
			)
		})
	}

	wg.Wait()

	// Run concurrent clients to consume the database
	// Check consistency in the database (no overbooking, etc.)
}
