package service

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// User Service
type User struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

// VerifyToken checks if user token is associated with a user and returns its id
func (u *User) VerifyToken(token string) int {

	var id int
	u.db.Table("users").Where("token = ?", token).Select("id").Row().Scan(&id)

	return id
}
