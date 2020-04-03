package routes

import (
	"github.com/gorilla/mux"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/api"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
)

// Register registers all routes with router
func Register(r *mux.Router, conf app.Config) {
	api := api.New(conf)

	r.HandleFunc("/resources", api.GetAllResources).Methods("GET")                                        //
	r.HandleFunc("/resource/{resourceID}/version/{versionID}", api.GetResourceByVersionID).Methods("GET") //
	r.HandleFunc("/categories", api.GetAllCategorieswithTags).Methods("GET")                              //
	r.HandleFunc("/resource/{resourceID}/rating", api.GetResourceRating).Methods("GET")                   //
	r.HandleFunc("/resource/{resourceID}/rating", api.UpdateResourceRating).Methods("PUT")                //
	r.HandleFunc("/oauth/redirect", api.GithubAuth).Methods("POST")                                       //

	r.HandleFunc("/resource/{id}", api.DeleteResourceHandler).Methods("DELETE")
	r.HandleFunc("/resource/yaml/{id}", api.GetResourceYAMLFile).Methods("GET")     //
	r.HandleFunc("/resource/readme/{id}", api.GetResourceReadmeFile).Methods("GET") //
	r.HandleFunc("/tags", api.GetAllTags).Methods("GET")                            //
	r.Path("/resources/{type}/{verified}").Queries("tags", "{tags}").HandlerFunc(api.GetAllFilteredResourcesByTag).Methods("GET")
	r.HandleFunc("/rating", api.AddRating).Methods("POST")   //
	r.HandleFunc("/upload", api.Upload).Methods("POST")      //
	r.HandleFunc("/stars", api.GetPrevStars).Methods("POST") //

	r.HandleFunc("/oauth/redirect", api.GithubAuth).Methods("POST")                       //
	r.HandleFunc("/resources/user/{id}", api.GetAllResourcesByUserHandler).Methods("GET") //
	r.HandleFunc("/resource/links/{id}", api.GetResourceLinksHandler).Methods("GET")      //
}
