package models

import (
	"database/sql"
	"fmt"
	"log"
)

// TaskTags represents many-many between Task and Tag models
type TaskTags struct {
	TaskID int `json:"taskID"`
	TagID  int `json:"tagID"`
}

// GetAllResourcesWithGivenTags queries for all resources with given tags
func GetAllResourcesWithGivenTags(resourceType string, verified string, tags []string) []Resource {
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
	var (
		resourceTagMap map[int][]string
		rows           *sql.Rows
		err            error
	)
	resourceTagMap = getResourceTagMap()
	rows, err = executeTagsQuery(tags, params, args)
	defer rows.Close()
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		resource := Resource{}
		err = rows.Scan(&resource.ID, &resource.Name, &resource.Type, &resource.Description, &resource.Downloads, &resource.Rating, &resource.Github, &resource.Verified)
		if err != nil {
			log.Println(err)
		}
		resource.Tags = resourceTagMap[resource.ID]
		matchTypeAndVerified(resourceType, verified, resource, &resources)
	}
	return resources
}

func executeTagsQuery(tags []string, params string, args []interface{}) (*sql.Rows, error) {
	var (
		rows         *sql.Rows
		err          error
		sqlStatement string
	)
	if len(tags) > 0 {
		sqlStatement = `
	SELECT DISTINCT T.ID,T.NAME,T.TYPE,T.DESCRIPTION,T.DOWNLOADS,T.RATING,T.GITHUB,T.VERIFIED
	FROM RESOURCE AS T JOIN RESOURCE_TAG AS TT ON (T.ID=TT.RESOURCE_ID) JOIN TAG
	AS TG ON (TG.ID=TT.TAG_ID AND TG.NAME in (` +
			params + `));`
		rows, err = DB.Query(sqlStatement, args...)
	} else {
		sqlStatement = `
	SELECT DISTINCT T.ID,T.NAME,T.TYPE,T.DESCRIPTION,T.DOWNLOADS,T.RATING,T.GITHUB,T.VERIFIED
	FROM RESOURCE T`
		rows, err = DB.Query(sqlStatement)
	}
	return rows, err
}

func matchTypeAndVerified(resourceType string, verified string, resource Resource, resources *[]Resource) {
	isVerified := getBoolString(resource.Verified)
	if resourceType != "all" && verified != "all" {
		if resourceType == resource.Type && isVerified == verified {
			*resources = append(*resources, resource)
		}
	} else if resourceType == "all" && verified != "all" {
		if isVerified == verified {
			*resources = append(*resources, resource)
		}
	} else if resourceType != "all" && verified == "all" {
		if resourceType == resource.Type {
			*resources = append(*resources, resource)
		}
	} else {
		*resources = append(*resources, resource)
	}
}

func getBoolString(p bool) string {
	if p == true {
		return "true"
	} else if p == false {
		return "false"
	}
	return "all"
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
