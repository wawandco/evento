package main

import (
	"fmt"
	"net/http"
	"os"

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

	fmt.Println("info: server running on :8080")
	err = http.ListenAndServe(":8080", srv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// Run concurrent clients to consume the database
	// Check consistency in the database (no overbooking, etc.)
}
