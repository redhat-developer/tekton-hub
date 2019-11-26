package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Pipelines-Marketplace/backend/pkg/models"
	"github.com/gorilla/mux"
)

// GetAllTasks writes json encoded tasks to ResponseWriter
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.GetAllTasks())
}

// GetTaskByID writes json encoded task to ResponseWriter
func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.GetTaskWithName(mux.Vars(r)["id"]))
}

// GetTaskFiles returns a compressed zip with task files
func GetTaskFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	GetCompressedFiles(mux.Vars(r)["name"])
	// Serve the created zip file
	http.ServeFile(w, r, "finalZipFile.zip")
}

// GetAllTags writes json encoded list of tags to Responsewriter
func GetAllTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.GetAllTags())
}

// GetAllFilteredTasksByTag writes json encoded list of filtered tasks to Responsewriter
func GetAllFilteredTasksByTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.GetAllTasksWithGivenTags(strings.Split(r.FormValue("tags"), "|")))
}

// GetAllFilteredTasksByCategory writes json encoded list of filtered tasks to Responsewriter
func GetAllFilteredTasksByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.GetAllTasksWithGivenCategory(strings.Split(r.FormValue("category"), "|")))
}

// GetTaskYAMLFile returns a compressed zip with task files
func GetTaskYAMLFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/file")
	files, err := ioutil.ReadDir("catalog" + "/" + mux.Vars(r)["name"])
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") {
			http.ServeFile(w, r, "catalog/"+mux.Vars(r)["name"]+"/"+f.Name())
			break
		}
	}
}

// GetTaskReadmeFile returns a compressed zip with task files
func GetTaskReadmeFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/file")
	files, err := ioutil.ReadDir("catalog" + "/" + mux.Vars(r)["name"])
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") {
			http.ServeFile(w, r, "catalog/"+mux.Vars(r)["name"]+"/"+f.Name())
			break
		}
	}
}

// DownloadFile returns a requested YAML file
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/file")
	requestedFileName := models.GetTaskNameFromID(mux.Vars(r)["id"])
	models.IncrementDownloads(mux.Vars(r)["id"])
	files, err := ioutil.ReadDir("catalog" + "/" + requestedFileName)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") {
			http.ServeFile(w, r, "catalog/"+requestedFileName+"/"+f.Name())
			break
		}
	}
}
