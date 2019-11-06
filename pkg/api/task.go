/*
This package handles the API requests related to tasks
*/
package api

import (
	"encoding/json"
	"net/http"
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allTasks())
}
