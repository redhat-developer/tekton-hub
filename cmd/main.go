package main

import (
	"log"
	"net/http"

	"github.com/Pipelines-Marketplace/backend/pkg/models"
	"github.com/Pipelines-Marketplace/backend/routers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	router := mux.NewRouter()
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := models.StartConnection(); err != nil {
		// log.Fatalln(err)
		panic(err)
	}
	// models.AddContentsToDB()
	routers.HandleRouters(router)
	http.ListenAndServe(":5000", router)
}
