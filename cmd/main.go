package main

import (
	"net/http"

	"github.com/Pipelines-Marketplace/backend/pkg/api"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	// Add new routers
	router.HandleFunc("/tasks", api.GetAllTasks).Methods("GET")
	http.ListenAndServe(":5000", router)
}
