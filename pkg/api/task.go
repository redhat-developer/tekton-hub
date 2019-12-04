package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Pipelines-Marketplace/backend/pkg/authentication"
	"github.com/Pipelines-Marketplace/backend/pkg/models"
	"github.com/Pipelines-Marketplace/backend/pkg/upload"
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
	taskID := mux.Vars(r)["id"]
	http.ServeFile(w, r, "tekton/"+taskID+".yaml")
}

// GetTaskReadmeFile returns a compressed zip with task files
func GetTaskReadmeFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/file")
	taskID := mux.Vars(r)["id"]
	readmeExists := models.DoesREADMEExist(taskID)
	if readmeExists {
		http.ServeFile(w, r, "readme/"+taskID+".md")
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"status": false, "message": "README file doesn't exist"})
}

// LoginHandler handles user authentication
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := &authentication.UserAuth{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
	}
	json.NewEncoder(w).Encode(authentication.Login(user))
}

// SignUpHandler registers a new user
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := &authentication.NewUser{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
	}
	json.NewEncoder(w).Encode(authentication.Signup(user))
}

// DownloadFile returns a requested YAML file
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/file")
	taskID := mux.Vars(r)["id"]
	models.IncrementDownloads(taskID)
	http.ServeFile(w, r, "tekton/"+taskID+".yaml")
}

// UpdateRating will add a new rating
func UpdateRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ratingRequestBody := AddRatingsRequest{}
	err := json.NewDecoder(r.Body).Decode(&ratingRequestBody)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(models.UpdateRating(ratingRequestBody.UserID, ratingRequestBody.TaskID, ratingRequestBody.Stars, ratingRequestBody.PrevStars))
}

// GetRatingDetails returns rating details of a task
func GetRatingDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.GetRatingDetialsByTaskID(mux.Vars(r)["id"]))
}

// AddRating add's a new rating
func AddRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ratingRequestBody := AddRatingsRequest{}
	err := json.NewDecoder(r.Body).Decode(&ratingRequestBody)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(models.AddRating(ratingRequestBody.UserID, ratingRequestBody.TaskID, ratingRequestBody.Stars, ratingRequestBody.PrevStars))
}

// Upload a new task/pipeline
func Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uploadRequestBody := upload.NewUploadRequestObject{}
	err := json.NewDecoder(r.Body).Decode(&uploadRequestBody)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(upload.NewUpload(uploadRequestBody.Name, uploadRequestBody.Description, uploadRequestBody.Type, uploadRequestBody.Tags, uploadRequestBody.Github, uploadRequestBody.UserID))
}
