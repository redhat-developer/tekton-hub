package models

import (
	"fmt"
	"log"
)

// TaskTags represents many-many between Task and Tag models
type TaskTags struct {
	TaskID int `json:"taskID"`
	TagID  int `json:"tagID"`
}

// GetAllResourcesWithGivenTags queries for all resources with given tags
func GetAllResourcesWithGivenTags(tags []string) []Resource {
	resources := []Resource{}
	args := make([]interface{}, len(tags))
	for index, value := range tags {
		args[index] = value
	}

	params := `$1`
	for index := 1; index <= len(tags); index++ {
		if index > 1 {
			params = params + fmt.Sprintf(",$%d", index)
		}
	}
	log.Println(params)
	var resourceTagMap map[int][]string
	resourceTagMap = make(map[int][]string)
	resourceTagMap = getResourceTagMap()
	sqlStatement := `
	SELECT DISTINCT T.ID,T.NAME,T.TYPE,T.DESCRIPTION,T.DOWNLOADS,T.RATING,T.GITHUB
	FROM RESOURCE AS T JOIN RESOURCE_TAG AS TT ON (T.ID=TT.RESOURCE_ID) JOIN TAG
	AS TG ON (TG.ID=TT.TAG_ID AND TG.NAME in (` +
		params + `));`
	log.Println(sqlStatement)
	rows, err := DB.Query(sqlStatement, args...)
	defer rows.Close()
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		resource := Resource{}
		err = rows.Scan(&resource.ID, &resource.Name, &resource.Type, &resource.Description, &resource.Downloads, &resource.Rating, &resource.Github)
		if err != nil {
			log.Println(err)
		}
		resource.Tags = resourceTagMap[resource.ID]
		resources = append(resources, resource)
	}
	return resources
}

func getResourceTagMap() map[int][]string {
	sqlStatement := `SELECT DISTINCT T.ID,TG.NAME FROM RESOURCE AS T JOIN RESOURCE_TAG AS TT ON (T.ID=TT.RESOURCE_ID) JOIN TAG AS TG ON (TG.ID=TT.TAG_ID);`
	rows, err := DB.Query(sqlStatement)
	// mapping task ID with tag names
	var taskTagMap map[int][]string
	taskTagMap = make(map[int][]string)
	for rows.Next() {
		var taskID int
		var tagName string
		err = rows.Scan(&taskID, &tagName)
		if err != nil {
			log.Println(err)
		}
		taskTagMap[taskID] = append(taskTagMap[taskID], tagName)
	}
	return taskTagMap
}
