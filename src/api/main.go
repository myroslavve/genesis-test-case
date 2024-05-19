package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Register handlers
	r.HandleFunc("/api/rate", RateHandler).Methods("GET")
	r.HandleFunc("/api/subscribe", SubscribeHandler).Methods("POST")

	// Start the server on port 8080
	log.Println("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %s\n", err.Error())
	}
}
