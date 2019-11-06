package models

type Task struct {
	Name        string
	id          int `gorm:"AUTO_INCREMENT"`
	Description string
	Downloads   int
	Rating      float64
}
