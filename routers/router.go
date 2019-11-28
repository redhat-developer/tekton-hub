package routers

import (
	"github.com/Pipelines-Marketplace/backend/pkg/api"
	"github.com/gorilla/mux"
)

// HandleRouters handle routes
func HandleRouters(router *mux.Router) {
	router.HandleFunc("/task/{id}", api.GetTaskByID).Methods("GET")
	router.HandleFunc("/task/{name}/files", api.GetTaskFiles).Methods("GET")
	router.HandleFunc("/task/{name}/yaml", api.GetTaskYAMLFile).Methods("GET")
	router.HandleFunc("/task/{name}/readme", api.GetTaskReadmeFile).Methods("GET")
	router.HandleFunc("/tags", api.GetAllTags).Methods("GET")
	router.Path("/tasks").Queries("tags", "{tags}").HandlerFunc(api.GetAllFilteredTasksByTag).Methods("GET")
	router.Path("/tasks").Queries("category", "{category}").HandlerFunc(api.GetAllFilteredTasksByCategory).Methods("GET")
	router.HandleFunc("/tasks", api.GetAllTasks).Methods("GET")
	router.HandleFunc("/download/{id}", api.DownloadFile).Methods("GET")
	router.HandleFunc("/rating/{id}", api.UpdateRating).Methods("POST")
	router.HandleFunc("/rating/{id}", api.GetRatingDetails).Methods("GET")
}
