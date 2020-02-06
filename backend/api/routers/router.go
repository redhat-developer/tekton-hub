package routers

import (
	"github.com/Pipelines-Marketplace/backend/pkg/api"
	"github.com/gorilla/mux"
)

// HandleRouters handle routes
func HandleRouters(router *mux.Router) {
	router.HandleFunc("/resource/{id}", api.GetResourceByID).Methods("GET") //
	router.HandleFunc("/resource/{id}", api.DeleteResourceHandler).Methods("DELETE")
	router.HandleFunc("/resource/yaml/{id}", api.GetResourceYAMLFile).Methods("GET")     //
	router.HandleFunc("/resource/readme/{id}", api.GetResourceReadmeFile).Methods("GET") //
	router.HandleFunc("/tags", api.GetAllTags).Methods("GET")                            //
	router.Path("/resources/{type}/{verified}").Queries("tags", "{tags}").HandlerFunc(api.GetAllFilteredResourcesByTag).Methods("GET")
	router.HandleFunc("/resources", api.GetAllResources).Methods("GET")    //
	router.HandleFunc("/rating", api.AddRating).Methods("POST")            //
	router.HandleFunc("/rating", api.UpdateRating).Methods("PUT")          //
	router.HandleFunc("/rating/{id}", api.GetRatingDetails).Methods("GET") //
	router.HandleFunc("/upload", api.Upload).Methods("POST")               //
	router.HandleFunc("/stars", api.GetPrevStars).Methods("POST")          //

	router.HandleFunc("/oauth/redirect", api.GithubAuth).Methods("POST")                       //
	router.HandleFunc("/resources/user/{id}", api.GetAllResourcesByUserHandler).Methods("GET") //
	router.HandleFunc("/resource/links/{id}", api.GetResourceLinksHandler).Methods("GET")      //
}
