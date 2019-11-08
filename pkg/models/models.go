package models

// Task is a database model representing task data
type Task struct {
	Name        string
	Description string
	Downloads   int
	Rating      float64
	YAML        string
}
