package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/routes"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/sync"
)

func main() {
	app, err := app.FromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: failed to initialise: %s", err)
		os.Exit(1)
	}
	defer app.Cleanup()

	db := app.DB()
	db.LogMode(true)

	s := sync.New(app)
	s.Init()

	// TODO(sthaha): start syncing all catalogs on startup
	//go func() {
	//catalog := model.Catalog{}
	//db.Model(&model.Catalog{}).First(&catalog)

	//job := model.SyncJob{Catalog: catalog, Status: "queued"}
	//db.Create(&job)
	//s.Sync(context.Background())
	//}()

	//HTTP
	router := mux.NewRouter()
	routes.Register(router, app, s)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{
			"X-Requested-With", "Content-Type", "Authorization",
		}),
		handlers.AllowedMethods([]string{
			"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE",
		}),
	)

	log := app.Logger()
	log.Infof("Listening on %s", app.Addr())
	log.Fatal(http.ListenAndServe(app.Addr(), cors(router)))
}
