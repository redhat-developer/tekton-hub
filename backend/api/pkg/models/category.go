package models

import "log"

// CategoryTag is a model representing categories associated with tags
type CategoryTag struct {
	Category string `json:"category"`
	Tag      string `json:"tag"`
}

// GetAllCategorieswithTags will query for all Categories
func GetAllCategorieswithTags() map[string][]string {

	categoryTagMap := make(map[string][]string)
	sqlStatement := `
	SELECT CATEGORY.NAME AS CATEGORY,TAG.NAME AS TAG FROM CATEGORY INNER JOIN TAG ON CATEGORY.ID=TAG.CATEGORY ORDER BY CATEGORY.ID;`
	rows, err := DB.Query(sqlStatement)
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
