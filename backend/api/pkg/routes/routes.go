package routes

import (
	"github.com/gorilla/mux"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/api"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
)

// Register registers all routes with router
func Register(r *mux.Router, conf app.Config) {
	api := api.New(conf)

	r.HandleFunc("/resource/{id}", api.GetResourceByID).Methods("GET") //
	r.HandleFunc("/resource/{id}", api.DeleteResourceHandler).Methods("DELETE")
	r.HandleFunc("/resource/yaml/{id}", api.GetResourceYAMLFile).Methods("GET")     //
	r.HandleFunc("/resource/readme/{id}", api.GetResourceReadmeFile).Methods("GET") //
	r.HandleFunc("/tags", api.GetAllTags).Methods("GET")                            //
	r.Path("/resources/{type}/{verified}").Queries("tags", "{tags}").HandlerFunc(api.GetAllFilteredResourcesByTag).Methods("GET")
	r.HandleFunc("/resources", api.GetAllResources).Methods("GET")    //
	r.HandleFunc("/rating", api.AddRating).Methods("POST")            //
	r.HandleFunc("/rating", api.UpdateRating).Methods("PUT")          //
	r.HandleFunc("/rating/{id}", api.GetRatingDetails).Methods("GET") //
	r.HandleFunc("/upload", api.Upload).Methods("POST")               //
	r.HandleFunc("/stars", api.GetPrevStars).Methods("POST")          //

	r.HandleFunc("/oauth/redirect", api.GithubAuth).Methods("POST")                       //
	r.HandleFunc("/resources/user/{id}", api.GetAllResourcesByUserHandler).Methods("GET") //
	r.HandleFunc("/resource/links/{id}", api.GetResourceLinksHandler).Methods("GET")      //
}
