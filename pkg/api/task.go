/*
This package handles the API requests related to tasks
*/
package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allTasks())
}

func GetTaskWithID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getTask(mux.Vars(r)["name"]))
}
