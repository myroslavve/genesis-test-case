package api

import (
	"github.com/gorilla/mux"
)

func InitializeRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/rate", RateHandler).Methods("GET")
	r.HandleFunc("/api/subscribe", SubscribeHandler).Methods("POST")
	return r
}
