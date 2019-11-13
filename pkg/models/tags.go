package models

import "log"

// Tag is a model representing tags associated with tasks
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetAllTags will query for all tags
func GetAllTags() []Tag {
	tag := []Tag{}
	if DB.HasTable([]Tag{}) {
		log.Println("Table exists!!")
	}
	err := DB.Find(&tag).RecordNotFound()
	log.Println(err)
	log.Println(tag)
	return tag
}
