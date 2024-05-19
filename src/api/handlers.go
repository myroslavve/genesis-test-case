package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/myroslavve/genesis-test-case/src/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Response struct {
	Message string `json:"message"`
}

var collection *mongo.Collection

func SetDB() {
	collection = db.GetCollection("subscriptions")
}

// Handler for the "/rate" endpoint
func RateHandler(w http.ResponseWriter, r *http.Request) {
	rate, err := FetchExchangeRate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err == nil {
		http.Error(w, "Email already subscribed", http.StatusConflict)
		return
	}

	_, err = collection.InsertOne(ctx, bson.M{"email": email})
	if err != nil {
		http.Error(w, "Failed to subscribe email", http.StatusInternalServerError)
		return
	}

	response := Response{Message: "Email subscribed"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
