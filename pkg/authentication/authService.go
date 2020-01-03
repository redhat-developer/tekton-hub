package authentication

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jWTToken struct {
	Token string
}

var authenticatedID int

var mySigningKey = []byte("supersecret")

// GenerateJWT a new JWT token
func GenerateJWT(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["id"] = userID
	claims["expiry"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenString, nil
}
