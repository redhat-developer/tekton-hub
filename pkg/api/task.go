/*
Package api handles the API requests related to tasks
*/
package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetAllTasks writes json encoded tasks to ResponseWriter
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allTasks())
}

// GetTaskWithID writes json encoded task to ResponseWriter
func GetTaskWithID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getTask(mux.Vars(r)["name"]))
}

// GetTaskFiles returns a compressed zip with task files
func GetTaskFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	GetCompressedFiles(mux.Vars(r)["name"])
	// Serve the created zip file
	http.ServeFile(w, r, "finalZipFile.zip")
}
