package models

import "log"

// Tag is a model representing tags associated with tasks
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetAllTags will query for all tags
func GetAllTags() []Tag {
	tags := []Tag{}
	sqlStatement := `
	SELECT * FROM TAG;`
	rows, err := DB.Query(sqlStatement)
	for rows.Next() {
		tag := Tag{}
		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			log.Println(err)
		}
		tags = append(tags, tag)
	}
	return tags
}
