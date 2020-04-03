package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/db/model"
	"go.uber.org/zap"
)

// User Service
type User struct {
	db  *gorm.DB
	log *zap.SugaredLogger
	gh  *app.GitHub
}

// GHUserDetails model represents user details
type GHUserDetails struct {
	UserName string `json:"login"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

// OAuthAuthorizeToken represents Authorise token
type OAuthAuthorizeToken struct {
	Token string `json:"token"`
}

// OAuthAccessToken represents Access token
type OAuthAccessToken struct {
	AccessToken string `json:"access_token"`
}

// OAuthResponse Api reponse
type OAuthResponse struct {
	Token string `json:"token"`
}

// VerifyToken checks if user token is associated with a user and returns its id
func (u *User) VerifyToken(token string) int {

	var id int
	u.db.Table("users").Where("token = ?", token).Select("id").Row().Scan(&id)

	return id
}

// Add insert user in database
func (u *User) Add(ud GHUserDetails) *model.User {

	user := &model.User{}
	u.db.Where("user_name = ?", ud.UserName).
		Assign(&model.User{Token: ud.Token}).
		FirstOrCreate(&model.User{
			UserName: ud.UserName,
			Name:     ud.Name,
			Email:    ud.Email,
			Token:    ud.Token,
		}).Scan(&user)

	return user
}

// GetOAuthURL return url for getting access token
func (u *User) GetOAuthURL(token string) string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		u.gh.OAuthClientID, u.gh.OAuthSecret, token)
}

// GetGitHubAccessToken fetch User GithubAccessToken
func (u *User) GetGitHubAccessToken(authToken OAuthAuthorizeToken) (string, error) {

	reqURL := u.GetOAuthURL(authToken.Token)
	u.log.Info(reqURL)

	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
	}
	req.Header.Set("accept", "application/json")

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
	}
	defer res.Body.Close()

	var oat OAuthAccessToken
	if err := json.NewDecoder(res.Body).Decode(&oat); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
	}
	u.log.Info("Access Token", oat.AccessToken)

	return oat.AccessToken, nil
}

// GetUserDetails fetch user details using GitHub Api
func (u *User) GetUserDetails(oat OAuthAccessToken) GHUserDetails {

	httpClient := http.Client{}
	reqURL := fmt.Sprintf("https://api.github.com/user")

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	req.Header.Set("Authorization", "token "+oat.AccessToken)
	if err != nil {
		u.log.Error(err)
	}
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("accept", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	userDetails := GHUserDetails{}
	if err := json.Unmarshal(body, &userDetails); err != nil {
		u.log.Error(err)
	}
	userDetails.Token = oat.AccessToken

	return userDetails
}

// GenerateJWT a new JWT token
func (u *User) GenerateJWT(user *model.User) (string, error) {

	jwtSigningKey := []byte(u.gh.JWTSigningKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["id"] = user.ID
	claims["name"] = user.Name
	claims["user_name"] = user.UserName

	tokenString, err := token.SignedString(jwtSigningKey)
	if err != nil {
		u.log.Info(err.Error)
		return "", err
	}
	return tokenString, nil
}
