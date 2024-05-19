package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Response struct for JSON responses
type Response struct {
    Message string `json:"message"`
}

// In-memory store for subscriptions
var emailSubscriptions = struct {
    sync.RWMutex
    emails map[string]bool
}{emails: make(map[string]bool)}

// Handler for the "/rate" endpoint
func RateHandler(w http.ResponseWriter, r *http.Request) {
    // Mock response for exchange rate, replace with actual API call
    rate := 27.5
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rate)
}

// Handler for the "/subscribe" endpoint
func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    if email == "" {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    emailSubscriptions.Lock()
    defer emailSubscriptions.Unlock()

    if _, exists := emailSubscriptions.emails[email]; exists {
        http.Error(w, "Email already subscribed", http.StatusConflict)
        return
    }

    emailSubscriptions.emails[email] = true

    response := Response{Message: "Email subscribed"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
