package main

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/myroslavve/genesis-test-case/src/api"
	"github.com/myroslavve/genesis-test-case/src/db"
	"github.com/myroslavve/genesis-test-case/src/services"
)

func main() {
	// Load the .env file
	envPath := filepath.Join(".", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the database
	db.InitDB()
	api.SetDB()
	services.InitEmailService()

	// Initialize routes
	r := api.InitializeRoutes()

	// Start the email service to send exchange rates every 24 hours
	go services.SendExchangeRates(24 * time.Hour)

	// Start the server on port 8080
	log.Println("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %s\n", err.Error())
	}
}
