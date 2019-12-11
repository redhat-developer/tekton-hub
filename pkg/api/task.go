package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	// Serve the file from github url
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
	json.NewEncoder(w).Encode("noreadme")
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

// GetPrevStars will return the previous rating
func GetPrevStars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	previousStarRequestBody := models.PrevStarRequest{}
	err := json.NewDecoder(r.Body).Decode(&previousStarRequestBody)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(models.GetUserRating(previousStarRequestBody.UserID, previousStarRequestBody.TaskID))

}

// OAuthAccessResponse represents access_token
type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

// Code will
type Code struct {
	Token string `json:"token"`
}

// GithubAuth handles OAuth by Github
func GithubAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	httpClient := http.Client{}
	var clientID string
	var clientSecret string
	clientID = "aac6161a58b4d7798f05"
	clientSecret = "138039bebc49df742f83f1c041126a0e53ec0e10"
	token := Code{}
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		log.Println(err)
	}
	// c, _ := ioutil.ReadAll(r.Body)
	log.Println("Code", token.Token)
	var code string
	code = token.Token
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
	log.Println(reqURL)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		println(os.Stdout, "could not create HTTP request: %v", err)
	}
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		println(os.Stdout, "could not send HTTP request: %v", err)
	}

	// Parse the request body into the `OAuthAccessResponse` struct
	var t OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
	}
	log.Println("Access Token", t.AccessToken)
	username, id := getUserDetails(t.AccessToken)
	log.Println(username, id)
	authToken, err := authentication.GenerateJWT(int(id))
	if err != nil {
		log.Println(err)
	}
	// Add user if doesn't exist
	sqlStatement := `SELECT EXISTS(SELECT 1 FROM USER_CREDENTIAL WHERE ID=$1)`
	var exists bool
	err = models.DB.QueryRow(sqlStatement, id).Scan(&exists)
	if err != nil {
		log.Println(err)
	}
	log.Println(exists)
	if !exists {
		sqlStatement := `INSERT INTO USER_CREDENTIAL(ID,USERNAME,FIRST_NAME) VALUES($1,$2,$3)`
		_, err := models.DB.Exec(sqlStatement, id, "github", "github")
		if err != nil {
			log.Println(err)
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"token": authToken, "user_id": int(id)})
}

func getUserDetails(accessToken string) (string, int) {
	httpClient := http.Client{}
	reqURL := fmt.Sprintf("https://api.github.com/user")
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	req.Header.Set("Authorization", "token "+accessToken)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Println(string(body))
	var userData map[string]interface{}
	if err := json.Unmarshal(body, &userData); err != nil {
		log.Println(err)
	}
	username := userData["login"].(string)
	id := userData["id"].(float64)
	log.Println(id)
	return string(username), int(id)
}
