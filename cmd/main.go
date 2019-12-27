package main

import (
	"log"
	"net/http"

	"github.com/Pipelines-Marketplace/backend/pkg/models"
	"github.com/Pipelines-Marketplace/backend/routers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// "f60343c9a243da51d4e247e7c228b7653d6d33e3"
	router := mux.NewRouter()
	if err := models.StartConnection(); err != nil {
		log.Fatalln(err)
	}
	// models.AddContentsToDB()
	// upload.GetAllYAMLFilesFromRepository("golang-build", "https://github.com/tektoncd/catalog")
	routers.HandleRouters(router)
	http.ListenAndServe(":5000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
}
