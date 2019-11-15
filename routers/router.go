package routers

import (
	"github.com/Pipelines-Marketplace/backend/pkg/api"
	"github.com/gorilla/mux"
)

// HandleRouters handle routes
func HandleRouters(router *mux.Router) {
	router.HandleFunc("/task/{id}", api.GetTaskByID).Methods("GET")
	router.HandleFunc("/task/{name}/files", api.GetTaskFiles).Methods("GET")
	router.HandleFunc("/tags", api.GetAllTags).Methods("GET")
	router.Path("/tasks").Queries("tags", "{tags}").HandlerFunc(api.GetAllFilteredTasks).Methods("GET")
	router.HandleFunc("/tasks", api.GetAllTasks).Methods("GET")
}
