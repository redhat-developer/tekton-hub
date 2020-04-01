package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/authentication"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/models"

	//"github.com/redhat-developer/tekton-hub/backend/api/pkg/models"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/polling"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/service"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/upload"
	"go.uber.org/zap"
)

const (
	errParseForm   string = "form-parse-error"
	errMissingKey  string = "key-not-found"
	errInvalidType string = "invalid-type"
)

type Api struct {
	app     app.Config
	Log     *zap.SugaredLogger
	service service.Service
}

func New(app app.Config) *Api {
	return &Api{
		app:     app,
		Log:     app.Logger().With("name", "api"),
		service: service.New(app),
	}
}

type ResponseError struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func (e *ResponseError) Error() string {
	return e.Detail
}

func intQueryVar(r *http.Request, key string, def int) (int, *ResponseError) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return def, nil
	}

	res, err := strconv.Atoi(value)
	if err != nil {
		return def, &ResponseError{
			Code:   errInvalidType,
			Detail: "query param " + key + " must be an int"}
	}

	return res, nil
}

func intPathVar(r *http.Request, key string) (int, *ResponseError) {
	value := mux.Vars(r)[key]

	res, err := strconv.Atoi(value)
	if err != nil {
		return 0, &ResponseError{
			Code:   errInvalidType,
			Detail: "Path param " + key + " must be an int"}
	}

	return res, nil
}

func invalidRequest(w http.ResponseWriter, status int, err *ResponseError) {
	type emptyList []interface{}

	res := struct {
		Data   emptyList       `json:"data"`
		Errors []ResponseError `json:"errors"`
	}{
		Data:   emptyList{},
		Errors: []ResponseError{*err},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func errorResponse(w http.ResponseWriter, err *ResponseError) {
	type emptyList []interface{}

	res := struct {
		Data   emptyList       `json:"data"`
		Errors []ResponseError `json:"errors"`
	}{
		Data:   emptyList{},
		Errors: []ResponseError{*err},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GetAllResources writes json encoded resources to ResponseWriter
func (api *Api) GetAllResources(w http.ResponseWriter, r *http.Request) {
	limit, err := intQueryVar(r, "limit", 100)
	if err != nil {
		invalidRequest(w, http.StatusBadRequest, err)
		return
	}

	resources, _ := api.service.Resource().All(service.Filter{Limit: limit})

	res := struct {
		Data   []service.ResourceDetail `json:"data"`
		Errors []ResponseError          `json:"errors"`
	}{
		Data:   resources,
		Errors: []ResponseError{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GetResourceByVersionID writes json encoded resources to ResponseWriter
func (api *Api) GetResourceByVersionID(w http.ResponseWriter, r *http.Request) {

	resourceID, err := intPathVar(r, "resourceID")
	if err != nil {
		invalidRequest(w, http.StatusBadRequest, err)
		return
	}
	versionID, err := intPathVar(r, "versionID")
	if err != nil {
		invalidRequest(w, http.StatusBadRequest, err)
		return
	}
	rv := service.ResourceVersion{
		ResourceID: resourceID,
		VersionID:  versionID,
	}

	resource, retErr := api.service.Resource().ByVersionID(rv)

	if retErr != nil {
		errorResponse(w, &ResponseError{Code: "invalid-input", Detail: retErr.Error()})
		return
	}

	res := struct {
		Data   service.ResourceVersionDetail `json:"data"`
		Errors []ResponseError               `json:"errors"`
	}{
		Data:   resource,
		Errors: []ResponseError{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GetAllCategorieswithTags writes json encoded list of categories to Responsewriter
func (api *Api) GetAllCategorieswithTags(w http.ResponseWriter, r *http.Request) {

	categories, _ := api.service.Category().All()

	res := struct {
		Data   []service.CategoryDetail `json:"data"`
		Errors []ResponseError          `json:"errors"`
	}{
		Data:   categories,
		Errors: []ResponseError{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GetResourceRating returns user's rating of a resource
func (api *Api) GetResourceRating(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	if token == "" {
		errorResponse(w, &ResponseError{Code: "invalid-header", Detail: "Token is missing in header"})
		return
	}
	resourceID, err1 := intPathVar(r, "resourceID")
	if err1 != nil {
		invalidRequest(w, http.StatusBadRequest, err1)
		return
	}

	userID := api.service.User().VerifyToken(token)
	if userID == 0 {
		errorResponse(w, &ResponseError{Code: "invalid-token", Detail: "User with associated token not found."})
		return
	}

	ids := service.UserResource{
		UserID:     userID,
		ResourceID: resourceID,
	}

	rating, _ := api.service.Rating().GetResourceRating(ids)

	res := struct {
		Data   service.RatingDetails `json:"data"`
		Errors []ResponseError       `json:"errors"`
	}{
		Data:   rating,
		Errors: []ResponseError{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// UpdateResourceRating will add/update a user's re rating
func (api *Api) UpdateResourceRating(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	if token == "" {
		errorResponse(w, &ResponseError{Code: "invalid-header", Detail: "Token is missing in header"})
		return
	}
	resourceID, err1 := intPathVar(r, "resourceID")
	if err1 != nil {
		invalidRequest(w, http.StatusBadRequest, err1)
		return
	}

	userID := api.service.User().VerifyToken(token)
	if userID == 0 {
		errorResponse(w, &ResponseError{Code: "invalid-token", Detail: "User with associated token not found."})
		return
	}

	ratingRequestBody := service.UpdateRatingDetails{UserID: uint(userID), ResourceID: uint(resourceID)}
	err := json.NewDecoder(r.Body).Decode(&ratingRequestBody)
	if err != nil {
		errorResponse(w, &ResponseError{Code: "invalid-body", Detail: err.Error()})
		return
	}

	if ratingRequestBody.ResourceRating > 5 {
		errorResponse(w, &ResponseError{Code: "invalid-rating", Detail: "Rating should be in range 1 to 5"})
		return
	}

	api.service.Rating().UpdateResourceRating(ratingRequestBody)

	type emptyList []interface{}
	res := struct {
		Data   emptyList       `json:"data"`
		Errors []ResponseError `json:"errors"`
	}{
		Data:   emptyList{},
		Errors: []ResponseError{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GetAllTags writes json encoded list of tags to Responsewriter
func (api *Api) GetAllTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.GetAllTags(api.app.DB()))
}

// GetAllFilteredResourcesByTag writes json encoded list of filtered tasks to Responsewriter
func (api *Api) GetAllFilteredResourcesByTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tags []string
	if r.FormValue("tags") != "" {
		tags = strings.Split(r.FormValue("tags"), "|")
	}
	json.NewEncoder(w).Encode(models.GetAllResourcesWithGivenTags(api.app.DB(), mux.Vars(r)["type"], mux.Vars(r)["verified"], tags))
}

// GetResourceYAMLFile returns a compressed zip with task files
func (api *Api) GetResourceYAMLFile(w http.ResponseWriter, r *http.Request) {
	resourceID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		api.Log.Error(err)
	}
	githubDetails := models.GetResourceGithubDetails(api.app.DB(), resourceID)
	gh := api.app.GitHub().Client
	desc, err := polling.GetFileContent(context.Background(), gh, githubDetails.Owner, githubDetails.RepositoryName, githubDetails.Path, nil)
	if err != nil {
		api.Log.Error(err)
		json.NewEncoder(w).Encode("noyaml")
		return
	}
	content, err := desc.GetContent()
	if err != nil {
		api.Log.Error(err)
		json.NewEncoder(w).Encode("noyaml")
		return
	}
	w.Write([]byte(content))
}

// GetResourceReadmeFile will return  a README file
func (api *Api) GetResourceReadmeFile(w http.ResponseWriter, r *http.Request) {
	resourceID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		api.Log.Error(err)
	}
	githubDetails := models.GetResourceGithubDetails(api.app.DB(), resourceID)
	if githubDetails.ReadmePath == "" {
		json.NewEncoder(w).Encode("noreadme")
		return
	}
	gh := api.app.GitHub().Client
	desc, err := polling.GetFileContent(context.Background(), gh, githubDetails.Owner, githubDetails.RepositoryName, githubDetails.ReadmePath, nil)
	if err != nil {
		api.Log.Error(err)
	}
	content, err := desc.GetContent()
	if err != nil {
		api.Log.Error(err)
	}
	w.Write([]byte(content))
}

// AddRating add's a new rating
func (api *Api) AddRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ratingRequestBody := AddRatingsRequest{}
	err := json.NewDecoder(r.Body).Decode(&ratingRequestBody)
	if err != nil {
		api.Log.Error(err)
	}
	json.NewEncoder(w).Encode(models.AddRating(api.app.DB(), ratingRequestBody.UserID, ratingRequestBody.ResourceID, ratingRequestBody.Stars, ratingRequestBody.PrevStars))
}

// Upload a new task/pipeline
func (api *Api) Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uploadRequestBody := upload.NewUploadRequestObject{}
	err := json.NewDecoder(r.Body).Decode(&uploadRequestBody)
	if err != nil {
		api.Log.Error(err)
	}
	uploader := upload.New(api.app)
	if uploadRequestBody.Type == "task" {
		json.NewEncoder(w).Encode(uploader.NewUpload(uploadRequestBody.Name, uploadRequestBody.Description, uploadRequestBody.Type, uploadRequestBody.Tags, uploadRequestBody.Github, uploadRequestBody.UserID))
	} else if uploadRequestBody.Type == "pipeline" {
		json.NewEncoder(w).Encode(uploader.NewUploadPipeline(uploadRequestBody.Name, uploadRequestBody.Description, uploadRequestBody.Type, uploadRequestBody.Tags, uploadRequestBody.Github, uploadRequestBody.UserID))
	}
}

// GetPrevStars will return the previous rating
func (api *Api) GetPrevStars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	previousStarRequestBody := models.PrevStarRequest{}
	err := json.NewDecoder(r.Body).Decode(&previousStarRequestBody)
	if err != nil {
		api.Log.Error(err)
	}
	json.NewEncoder(w).Encode(models.GetUserRating(api.app.DB(), previousStarRequestBody.UserID, previousStarRequestBody.ResourceID))

}

func ghOAuthURLForCode(code string) string {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	return fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		clientID, clientSecret, code)
}

// GithubAuth handles OAuth by Github
func (api *Api) GithubAuth(w http.ResponseWriter, r *http.Request) {

	token := Code{}
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		api.Log.Error(err)
	}
	api.Log.Info("Code", token.Token)

	reqURL := ghOAuthURLForCode(token.Token)
	api.Log.Info(reqURL)

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
	api.Log.Info("Access Token", t.AccessToken)
	username, id := api.getUserDetails(t.AccessToken)
	api.Log.Info(username, id)
	authToken, err := authentication.GenerateJWT(int(id))
	if err != nil {
		api.Log.Error(err)
	}

	// Add user if doesn't exist
	sqlStatement := `SELECT EXISTS(SELECT 1 FROM USER_CREDENTIAL WHERE ID=$1)`
	var exists bool
	db := api.app.DB().DB()
	err = db.QueryRow(sqlStatement, id).Scan(&exists)
	if err != nil {
		api.Log.Error(err)
	}
	api.Log.Info(exists)

	if !exists {
		sqlStatement := `INSERT INTO USER_CREDENTIAL(ID,USERNAME,FIRST_NAME,TOKEN) VALUES($1,$2,$3,$4)`
		_, err := db.Exec(sqlStatement, id, "github", "github", t.AccessToken)
		if err != nil {
			api.Log.Error(err)
		}
	} else {
		// Update token if user exists
		sqlStatement = `UPDATE USER_CREDENTIAL SET TOKEN=$2 WHERE ID=$1`
		_, err = db.Exec(sqlStatement, id, t.AccessToken)
		if err != nil {
			api.Log.Error(err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"token": authToken, "user_id": int(id)})
}

func (api *Api) getUserDetails(accessToken string) (string, int) {
	httpClient := http.Client{}
	reqURL := fmt.Sprintf("https://api.github.com/user")
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	req.Header.Set("Authorization", "token "+accessToken)
	if err != nil {
		api.Log.Error(err)
	}
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		api.Log.Error(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	api.Log.Info(string(body))
	var userData map[string]interface{}
	if err := json.Unmarshal(body, &userData); err != nil {
		api.Log.Error(err)
	}
	username := userData["login"].(string)
	id := userData["id"].(float64)
	return string(username), int(id)
}

// GetAllResourcesByUserHandler will return all tasks uploaded by user
func (api *Api) GetAllResourcesByUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"status": false, "message": "Invalid User ID"})
	}
	json.NewEncoder(w).Encode(models.GetAllResourcesByUser(api.app.DB(), userID))
}

// DeleteResourceHandler handles resource deletion
func (api *Api) DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resourceID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		api.Log.Error(err)
	}
	err = models.DeleteResource(api.app.DB(), resourceID)
	if err != nil {
		api.Log.Error(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": true, "message": err})
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"status": true, "message": "Successfully Deleted"})
}

// GetResourceLinksHandler will return raw github links
func (api *Api) GetResourceLinksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resourceID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		api.Log.Error(err)
	}
	links := models.GetResourceRawLinks(api.app.DB(), resourceID)
	json.NewEncoder(w).Encode(links)
}
