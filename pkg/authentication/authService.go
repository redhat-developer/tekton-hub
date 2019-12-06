package authentication

import (
	"log"
	"time"

	"github.com/Pipelines-Marketplace/backend/pkg/models"
	"github.com/dgrijalva/jwt-go"
)

// UserAuth represents user model for authentication
type UserAuth struct {
	Username string
	Password string
}

type jWTToken struct {
	Token string
}

var authenticatedID int

// Login users
func Login(user *UserAuth) interface{} {
	// Check if username and password exists
	if !authenticate(user) {
		return map[string]interface{}{"status": false, "message": "Invalid credentials"}
	}
	token, err := GenerateJWT(user)
	if err != nil {
		log.Println(err)
	}
	return map[string]interface{}{"token": token, "user_id": authenticatedID}
}

var mySigningKey = []byte("supersecret")

// GenerateJWT a new JWT token
func GenerateJWT(user *UserAuth) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = user.Username
	claims["expiry"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenString, nil
}

// Authenticate user
func authenticate(user *UserAuth) bool {
	var id int
	sqlStatement := `SELECT ID FROM USER_AUTH WHERE USERNAME=$1 AND PASSWORD=crypt($2, PASSWORD)`
	err := models.DB.QueryRow(sqlStatement, user.Username, user.Password).Scan(&id)
	if err != nil {
		return false
	}
	authenticatedID = id
	return true
}
