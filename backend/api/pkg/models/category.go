package models

import (
	"log"
)

// Category is a model representing categories associated with tags
type Category struct {
	ID   int    `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:100;not null;unique" json:"name"`
}

// CategoryTag Model
type CategoryTag struct {
	Category string `json:"category"`
	Tag      string `json:"tag"`
}

// GetAllCategorieswithTags will query for all Categories with their associated tags
func GetAllCategorieswithTags() map[string][]string {
	categoryTagMap := make(map[string][]string)
	rows, err := GDB.Table("category").Select("category.name as category, tag.name as tag").Joins("inner join tag on tag.category_id = category.id").Rows()
	for rows.Next() {
		category := CategoryTag{}
		err = rows.Scan(&category.Category, &category.Tag)
		if err != nil {
			log.Println(err)
		}
		categoryTagMap[category.Category] = append(categoryTagMap[category.Category], category.Tag)
	}
	return categoryTagMap
}
