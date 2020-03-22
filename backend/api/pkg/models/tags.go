package models

import "log"

// Tag is a model representing tags associated with tasks
type Tag struct {
	ID         int    `gorm:"primary_key;auto_increment" json:"id"`
	Name       string `gorm:"not null;unique" json:"name"`
	CategoryID int    `json:"category_id"`
}

// GetAllTags will query for all tags
func GetAllTags() []Tag {
	tags := []Tag{}
	sqlStatement := `
	SELECT * FROM TAG;`
	rows, err := DB.Query(sqlStatement)
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
func AddTag(tag string) (int, error) {
	var newTagID int
	categoryID := 8
	sqlStatement := `INSERT INTO TAG(NAME,CATEGORY_ID) VALUES($1, $2) RETURNING ID`
	err := DB.QueryRow(sqlStatement, tag, categoryID).Scan(&newTagID)
	if err != nil {
		return 0, err
	}
	return newTagID, nil
}
