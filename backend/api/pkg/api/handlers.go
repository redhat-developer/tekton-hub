package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"

	//"github.com/redhat-developer/tekton-hub/backend/api/pkg/models"

	"github.com/redhat-developer/tekton-hub/backend/api/pkg/service"
	"go.uber.org/zap"
)

const (
	errParseForm   string = "form-parse-error"
	errMissingKey  string = "key-not-found"
	errInvalidType string = "invalid-type"
	errAuthFailure string = "auth-failed"
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
		err := &ResponseError{Code: "invalid-header", Detail: "Token is missing in header"}
		invalidRequest(w, http.StatusBadRequest, err)
		return
	}
	resourceID, err1 := intPathVar(r, "resourceID")
	if err1 != nil {
		invalidRequest(w, http.StatusBadRequest, err1)
		return
	}

	userID := api.service.User().VerifyToken(token)
	if userID == 0 {
		err := &ResponseError{Code: "invalid-header", Detail: "User with associated token not found."}
		invalidRequest(w, http.StatusBadRequest, err)
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
		err := &ResponseError{Code: "invalid-header", Detail: "Token is missing in header"}
		invalidRequest(w, http.StatusBadRequest, err)
		return
	}
	resourceID, err := intPathVar(r, "resourceID")
	if err != nil {
		invalidRequest(w, http.StatusBadRequest, err)
		return
	}

	userID := api.service.User().VerifyToken(token)
	if userID == 0 {
		err := &ResponseError{Code: "invalid-header", Detail: "User with associated token not found."}
		invalidRequest(w, http.StatusBadRequest, err)
		return
	}

	ratingRequestBody := service.UpdateRatingDetails{UserID: uint(userID), ResourceID: uint(resourceID)}
	jsonErr := json.NewDecoder(r.Body).Decode(&ratingRequestBody)
	if jsonErr != nil {
		err := &ResponseError{Code: "invalid-body", Detail: jsonErr.Error()}
		invalidRequest(w, http.StatusBadRequest, err)
		return
	}

	if ratingRequestBody.ResourceRating > 5 {
		err := &ResponseError{Code: "invalid-body", Detail: "Rating should be in range 1 to 5"}
		invalidRequest(w, http.StatusBadRequest, err)
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

// GithubAuth handles OAuth by Github
func (api *Api) GithubAuth(w http.ResponseWriter, r *http.Request) {

	token := service.OAuthAuthorizeToken{}
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		err := &ResponseError{Code: "invalid-body", Detail: err.Error()}
		invalidRequest(w, http.StatusBadRequest, err)
		return
	}
	api.Log.Info("OAuthAuthorizeToken - ", token.Token)

	accessToken, err := api.service.User().GetGitHubAccessToken(token)
	if err != nil {
		err := &ResponseError{Code: "invalid-header", Detail: "Invalid Token"}
		invalidRequest(w, http.StatusBadRequest, err)
		return
	}

	userDetails := api.service.User().GetUserDetails(service.OAuthAccessToken{AccessToken: accessToken})

	user := api.service.User().Add(userDetails)

	resToken, _ := api.service.User().GenerateJWT(user)

	res := struct {
		Data   service.OAuthResponse `json:"data"`
		Errors []ResponseError       `json:"errors"`
	}{
		Data:   service.OAuthResponse{Token: resToken},
		Errors: []ResponseError{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
