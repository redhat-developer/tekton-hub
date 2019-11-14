package main

import (
	"log"
	"net/http"

	"github.com/Pipelines-Marketplace/backend/pkg/models"
	"github.com/Pipelines-Marketplace/backend/routers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	// Add new routers
	if err := models.StartConnection(); err != nil {
		log.Fatalln(err)
	}
	// models.AddContentsToDB()
	routers.HandleRouters(router)
	http.ListenAndServe(":5000", router)
}
