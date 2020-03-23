package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

// Tag is a model representing tags associated with tasks
type Tag struct {
	ID         int    `gorm:"primary_key;auto_increment" json:"id"`
	Name       string `gorm:"not null;unique" json:"name"`
	CategoryID int    `json:"category_id"`
}

// GetAllTags will query for all tags
func GetAllTags(db *gorm.DB) []Tag {
	tags := []Tag{}
	sqlStatement := `
	SELECT * FROM TAG;`
	rows, err := db.DB().Query(sqlStatement)
	for rows.Next() {
		tag := Tag{}
		err = rows.Scan(&tag.ID, &tag.Name, &tag.CategoryID)
		if err != nil {
			log.Println(err)
		}
		tags = append(tags, tag)
	}
	return tags
}

// AddTag will add a new tag
func AddTag(db *gorm.DB, tag string) (int, error) {
	var newTagID int
	sqlStatement := `INSERT INTO TAG(NAME) VALUES($1) RETURNING ID`
	err := db.DB().QueryRow(sqlStatement, tag).Scan(&newTagID)
	if err != nil {
		return 0, err
	}
	return newTagID, nil
}
