package main

import (
	"log"
	"net/http"

	"github.com/redhat-developer/tekton-hub/backend/api/pkg/models"
	"github.com/redhat-developer/tekton-hub/backend/api/routers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	if err := models.StartConnection(); err != nil {
		log.Fatalln(err)
	}
	// models.AddResourcesFromCatalog("tektoncd", "catalog")
	routers.HandleRouters(router)
	http.ListenAndServe(":5000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}), handlers.AllowedOrigins([]string{"*"}))(router))
}
