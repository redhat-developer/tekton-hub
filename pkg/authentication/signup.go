package authentication

import (
	"log"

	"github.com/Pipelines-Marketplace/backend/pkg/models"
)

// NewUser represents a user requesting for signup
type NewUser struct {
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Signup will register a new user
func Signup(newUser *NewUser) interface{} {
	sqlStatement := `INSERT INTO USER_CREDENTIAL(USERNAME,FIRST_NAME,LAST_NAME,EMAIL) VALUES($1,$2,$3,$4)`
	_, err := models.DB.Exec(sqlStatement, newUser.Username, newUser.FirstName, newUser.LastName, newUser.Email)
	if err != nil {
		log.Println(err)
		return map[string]interface{}{"status": "false", "message": "Username already exists"}
	}
	var id int
	sqlStatement = `SELECT ID FROM USER_CREDENTIAL WHERE USERNAME=$1`
	err = models.DB.QueryRow(sqlStatement, newUser.Username).Scan(&id)
	if err != nil {
		log.Println(err)
	}
	addCredential(newUser.Username, newUser.Password, id)
	return map[string]interface{}{"status": "true", "message": "You've been registered"}
}

func addCredential(username string, password string, userID int) {
	sqlStatement := `INSERT INTO USER_AUTH(ID,USERNAME,PASSWORD) VALUES($1,$2,crypt($3, gen_salt('bf')))`
	_, err := models.DB.Exec(sqlStatement, userID, username, password)
	if err != nil {
		log.Println(err)
	}
}
