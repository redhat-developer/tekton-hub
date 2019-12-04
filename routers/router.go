package routers

import (
	"github.com/Pipelines-Marketplace/backend/pkg/api"
	"github.com/gorilla/mux"
)

// HandleRouters handle routes
func HandleRouters(router *mux.Router) {
	router.HandleFunc("/task/{id}", api.GetTaskByID).Methods("GET")
	router.HandleFunc("/task/{name}/files", api.GetTaskFiles).Methods("GET")
	router.HandleFunc("/task/{id}/yaml", api.GetTaskYAMLFile).Methods("GET")
	router.HandleFunc("/task/{id}/readme", api.GetTaskReadmeFile).Methods("GET")
	router.HandleFunc("/tags", api.GetAllTags).Methods("GET")
	router.Path("/tasks").Queries("tags", "{tags}").HandlerFunc(api.GetAllFilteredTasksByTag).Methods("GET")
	router.Path("/tasks").Queries("category", "{category}").HandlerFunc(api.GetAllFilteredTasksByCategory).Methods("GET")
	router.HandleFunc("/tasks", api.GetAllTasks).Methods("GET")
	router.HandleFunc("/login", api.LoginHandler).Methods("POST")
	router.HandleFunc("/signup", api.SignUpHandler).Methods("POST")
	router.HandleFunc("/rating", api.AddRating).Methods("POST")
	router.HandleFunc("/rating", api.UpdateRating).Methods("PUT")
	router.HandleFunc("/rating/{id}", api.GetRatingDetails).Methods("GET")
	router.HandleFunc("/download/{id}", api.DownloadFile).Methods("POST")
	router.HandleFunc("/upload", api.Upload).Methods("POST")
	router.HandleFunc("/stars", api.GetPrevStars).Methods("POST")
}
