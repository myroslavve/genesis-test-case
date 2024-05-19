package main

import (
	"log"
	"net/http"

	"github.com/myroslavve/genesis-test-case/src/api"
	"github.com/myroslavve/genesis-test-case/src/db"
)

func main() {
	// Initialize the database
	db.InitDB()
	api.SetDB()

	// Initialize routes
	r := api.InitializeRoutes()

	// Start the server on port 8080
	log.Println("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %s\n", err.Error())
	}
}
