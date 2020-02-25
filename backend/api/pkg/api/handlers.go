package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/authentication"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/models"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/polling"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/upload"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/utility"
)

// GetAllResources writes json encoded resources to ResponseWriter
func GetAllResources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.GetAllResources())
}

// GetResourceByID writes json encoded resources to ResponseWriter
func GetResourceByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resourceID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"status": false, "message": "Invalid User ID"})
	}
	json.NewEncoder(w).Encode(models.GetResourceByID(resourceID))
}

// GetAllTags writes json encoded list of tags to Responsewriter
func GetAllTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.GetAllTags())
}

// GetAllFilteredResourcesByTag writes json encoded list of filtered tasks to Responsewriter
func GetAllFilteredResourcesByTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tags []string
	if r.FormValue("tags") != "" {
		tags = strings.Split(r.FormValue("tags"), "|")
	}
	json.NewEncoder(w).Encode(models.GetAllResourcesWithGivenTags(mux.Vars(r)["type"], mux.Vars(r)["verified"], tags))
}

// GetResourceYAMLFile returns a compressed zip with task files
func GetResourceYAMLFile(w http.ResponseWriter, r *http.Request) {
	resourceID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
	}
	githubDetails := models.GetResourceGithubDetails(resourceID)
	desc, err := polling.GetFileContent(utility.Ctx, utility.Client, githubDetails.Owner, githubDetails.RepositoryName, githubDetails.Path, nil)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode("noyaml")
		return
	}
	content, err := desc.GetContent()
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode("noyaml")
		return
	}
	w.Write([]byte(content))
}

// GetResourceReadmeFile will return  a README file
func GetResourceReadmeFile(w http.ResponseWriter, r *http.Request) {
	resourceID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
	}
	githubDetails := models.GetResourceGithubDetails(resourceID)
	if githubDetails.ReadmePath == "" {
		json.NewEncoder(w).Encode("noreadme")
		return
	}
	desc, err := polling.GetFileContent(utility.Ctx, utility.Client, githubDetails.Owner, githubDetails.RepositoryName, githubDetails.ReadmePath, nil)
	if err != nil {
		log.Println(err)
	}
	content, err := desc.GetContent()
	if err != nil {
		log.Println(err)
	}
	w.Write([]byte(content))
}

// UpdateRating will add a new rating
func UpdateRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ratingRequestBody := AddRatingsRequest{}
	err := json.NewDecoder(r.Body).Decode(&ratingRequestBody)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(models.UpdateRating(ratingRequestBody.UserID, ratingRequestBody.ResourceID, ratingRequestBody.Stars, ratingRequestBody.PrevStars))
}

// GetRatingDetails returns rating details of a task
func GetRatingDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resourceID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(models.GetRatingDetialsByResourceID(resourceID))
}

// AddRating add's a new rating
func AddRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ratingRequestBody := AddRatingsRequest{}
	err := json.NewDecoder(r.Body).Decode(&ratingRequestBody)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(models.AddRating(ratingRequestBody.UserID, ratingRequestBody.ResourceID, ratingRequestBody.Stars, ratingRequestBody.PrevStars))
}

// Upload a new task/pipeline
func Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uploadRequestBody := upload.NewUploadRequestObject{}
	err := json.NewDecoder(r.Body).Decode(&uploadRequestBody)
	if err != nil {
		log.Println(err)
	}
	if uploadRequestBody.Type == "task" {
		json.NewEncoder(w).Encode(upload.NewUpload(uploadRequestBody.Name, uploadRequestBody.Description, uploadRequestBody.Type, uploadRequestBody.Tags, uploadRequestBody.Github, uploadRequestBody.UserID))
	} else if uploadRequestBody.Type == "pipeline" {
		json.NewEncoder(w).Encode(upload.NewUploadPipeline(uploadRequestBody.Name, uploadRequestBody.Description, uploadRequestBody.Type, uploadRequestBody.Tags, uploadRequestBody.Github, uploadRequestBody.UserID))
	}
}

// GetPrevStars will return the previous rating
func GetPrevStars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	previousStarRequestBody := models.PrevStarRequest{}
	err := json.NewDecoder(r.Body).Decode(&previousStarRequestBody)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(models.GetUserRating(previousStarRequestBody.UserID, previousStarRequestBody.ResourceID))

}

func ghOAuthURLForCode(code string) string {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	return fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		clientID, clientSecret, code)
}

// GithubAuth handles OAuth by Github
func GithubAuth(w http.ResponseWriter, r *http.Request) {

	token := Code{}
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		log.Println(err)
	}
	log.Println("Code", token.Token)

	reqURL := ghOAuthURLForCode(token.Token)
	log.Println(reqURL)

	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
	}
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	httpClient := http.Client{}
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
		sqlStatement := `INSERT INTO USER_CREDENTIAL(ID,USERNAME,FIRST_NAME,TOKEN) VALUES($1,$2,$3,$4)`
		_, err := models.DB.Exec(sqlStatement, id, "github", "github", t.AccessToken)
		if err != nil {
			log.Println(err)
		}
	} else {
		// Update token if user exists
		sqlStatement = `UPDATE USER_CREDENTIAL SET TOKEN=$2 WHERE ID=$1`
		_, err = models.DB.Exec(sqlStatement, id, t.AccessToken)
		if err != nil {
			log.Println(err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
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
	return string(username), int(id)
}

// GetAllResourcesByUserHandler will return all tasks uploaded by user
func GetAllResourcesByUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"status": false, "message": "Invalid User ID"})
	}
	json.NewEncoder(w).Encode(models.GetAllResourcesByUser(userID))
}

// DeleteResourceHandler handles resource deletion
func DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resourceID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
	}
	err = models.DeleteResource(resourceID)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": true, "message": err})
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"status": true, "message": "Successfully Deleted"})
}

// GetResourceLinksHandler will return raw github links
func GetResourceLinksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resourceID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
	}
	links := models.GetResourceRawLinks(resourceID)
	json.NewEncoder(w).Encode(links)
}
