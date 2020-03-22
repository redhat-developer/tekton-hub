package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/models"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/routes"
)

func main() {
	app, err := app.FromEnv("api")
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: failed to initialise: %s", err)
		os.Exit(1)
	}

	log := app.Logger()
	defer log.Sync()

	if err := models.Connect(app); err != nil {
		log.Fatalf("db connection failed: %s", err)
	}
	defer models.DB.Close()

	router := mux.NewRouter()
	// models.AddResourcesFromCatalog("tektoncd", "catalog")
	routes.Register(router, app)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{
			"X-Requested-With", "Content-Type", "Authorization",
		}),
		handlers.AllowedMethods([]string{
			"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE",
		}),
	)

	log.Infof("Listening on %s", app.Addr())
	log.Fatal(http.ListenAndServe(app.Addr(), cors(router)))
}
