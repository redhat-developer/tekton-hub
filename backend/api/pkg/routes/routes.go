package routes

import (
	"github.com/gorilla/mux"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/api"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/sync"
)

// Register registers all routes with router
func Register(r *mux.Router, conf app.Config, sync *sync.Service) {
	api := api.New(conf, sync)

	r.HandleFunc("/resources", api.GetAllResources).Methods("GET")                          //
	r.HandleFunc("/resource/{resourceID}/versions", api.GetResourceVersions).Methods("GET") //
	r.HandleFunc("/categories", api.GetAllCategorieswithTags).Methods("GET")                //
	r.HandleFunc("/resource/{resourceID}/rating", api.GetResourceRating).Methods("GET")     //
	r.HandleFunc("/resource/{resourceID}/rating", api.UpdateResourceRating).Methods("PUT")  //
	r.HandleFunc("/oauth/redirect", api.GithubAuth).Methods("POST")                         //
	r.HandleFunc("/resources/sync", api.SyncResources).Methods("POST")                      //

}
