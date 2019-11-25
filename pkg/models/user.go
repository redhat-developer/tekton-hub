package models

// User represents User model in database
type User struct {
	ID         int    `json:"id"`
	FirstName  string `json:"username"`
	SecondName string `json:"password"`
	EMAIL      string `json:"email"`
}
