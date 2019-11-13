package models

import "log"

// TaskTags represents many-many between Task and Tag models
type TaskTags struct {
	TaskID int
	TagID  int
}

// T ds
type T struct {
	ID          int
	Name        string
	Description string
	Downloads   int
	Ratings     int
}

// GetAllTasksWithGivenTags queries for all tasks with given tags
func GetAllTasksWithGivenTags(tags []string) []string {
	tasks := []T{}
	//DB.Where().Find()
	log.Println(tasks)
	return tags
}
