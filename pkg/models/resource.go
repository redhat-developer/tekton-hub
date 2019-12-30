package models

import (
	"log"
	"strconv"
)

// Resource is a database model representing task and pipeline
type Resource struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Downloads   int      `json:"downloads"`
	Rating      float64  `json:"rating"`
	Github      string   `json:"github"`
	Tags        []string `json:"tags"`
	Verified    bool     `json:"verified"`
}

// AddCatalogResource is called to add resource from catalog
func AddCatalogResource(resource *Resource) (int, error) {
	sqlStatement := `
	INSERT INTO RESOURCE (NAME,DOWNLOADS,RATING,GITHUB) 
	VALUES ($1, $2, $3, $4) RETURNING ID`
	var resourceID int
	err := DB.QueryRow(sqlStatement, resource.Name, resource.Downloads, resource.Rating, resource.Github).Scan(&resourceID)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return resourceID, nil
}

// AddResource will add a new resource
func AddResource(resource *Resource, userID int, owner string, respositoryName string, path string) (int, error) {
	var resourceID int
	sqlStatement := `
	INSERT INTO RESOURCE (NAME,DESCRIPTION,DOWNLOADS,RATING,GITHUB)
	VALUES ($1, $2, $3, $4, $5) RETURNING ID`
	err := DB.QueryRow(sqlStatement, resource.Name, resource.Description, resource.Downloads, resource.Rating, resource.Github).Scan(&resourceID)
	if err != nil {
		return 0, err
	}
	// Add Tags separately
	if len(resource.Tags) > 0 {
		for _, tag := range resource.Tags {
			tagExists := true
			// Use existing tags if already exists
			var tagID int
			sqlStatement = `SELECT ID FROM TAG WHERE NAME=$1`
			err := DB.QueryRow(sqlStatement, tag).Scan(&tagID)
			if err != nil {
				tagExists = false
				log.Println(err)
			}
			// If tag already exists
			if tagExists {
				addResourceTag(resourceID, tagID)
			} else {
				var newTagID int
				newTagID, err = AddTag(tag)
				if err != nil {
					log.Println(err)
				}
				addResourceTag(resourceID, newTagID)
			}
		}
	}
	addGithubDetails(resourceID, owner, respositoryName, path)
	return resourceID, addUserResource(userID, resourceID)
}

func addResourceTag(resourceID int, tagID int) {
	sqlStatement := `INSERT INTO RESOURCE_TAG(RESOURCE_ID,TAG_ID) VALUES($1,$2)`
	_, err := DB.Exec(sqlStatement, resourceID, tagID)
	if err != nil {
		log.Println(err)
	}
}

func addGithubDetails(resourceID int, owner string, respositoryName string, path string) {
	sqlStatement := `INSERT INTO GITHUB_DETAIL(RESOURCE_ID,OWNER,REPOSITORY_NAME,PATH) VALUES($1,$2,$3,$4)`
	_, err := DB.Exec(sqlStatement, resourceID, owner, respositoryName, path)
	if err != nil {
		log.Println(err)
	}
}

func updateGithubYAMLDetails(resourceID int, path string) {
	sqlStatement := `UPDATE GITHUB_DETAIL SET PATH=$1 WHERE RESOURCE_ID=$2`
	_, err := DB.Exec(sqlStatement, path, resourceID)
	if err != nil {
		log.Println(err)
	}
}

func updateGithubREADMEDetails(resourceID int, path string) {
	sqlStatement := `UPDATE GITHUB_DETAIL SET README_PATH=$1 WHERE RESOURCE_ID=$2`
	_, err := DB.Exec(sqlStatement, path, resourceID)
	if err != nil {
		log.Println(err)
	}
}

func addUserResource(userID int, resourceID int) error {
	sqlStatement := `INSERT INTO USER_RESOURCE(RESOURCE_ID,USER_ID) VALUES($1,$2)`
	_, err := DB.Exec(sqlStatement, resourceID, userID)
	if err != nil {
		return err
	}
	return nil
}

// CheckSameResourceUpload will checkif the user submitted the same resource again
func CheckSameResourceUpload(userID int, name string) bool {
	sqlStatement := `SELECT T.NAME FROM RESOURCE T JOIN USER_RESOURCE U ON T.ID=U.RESOURCE_ID WHERE U.USER_ID=$1`
	rows, err := DB.Query(sqlStatement, userID)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var resourceName string
		err := rows.Scan(&resourceName)
		if err != nil {
			log.Println(err)
		}
		if resourceName == name {
			return true
		}
	}
	return false
}

// GetAllResources will return all the tasks
func GetAllResources() []Resource {
	resources := []Resource{}
	sqlStatement := `
	SELECT * FROM RESOURCE ORDER BY ID`
	rows, err := DB.Query(sqlStatement)
	defer rows.Close()
	for rows.Next() {
		resource := Resource{}
		err = rows.Scan(&resource.ID, &resource.Name, &resource.Description, &resource.Downloads, &resource.Rating, &resource.Github, &resource.Type, &resource.Verified)
		if err != nil {
			log.Println(err)
		}
		resources = append(resources, resource)
	}
	resourceIndexMap := make(map[int]int)
	sqlStatement = `SELECT ID FROM RESOURCE ORDER BY ID`
	rows, err = DB.Query(sqlStatement)
	if err != nil {
		log.Println(err)
	}
	resourceIndex := 0
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		resourceIndexMap[id] = resourceIndex
		resourceIndex = resourceIndex + 1
	}

	sqlStatement = `SELECT R.ID,TG.NAME FROM TAG TG JOIN RESOURCE_TAG TT ON TT.TAG_ID=TG.ID JOIN RESOURCE R ON R.ID=TT.RESOURCE_ID`
	rows, err = DB.Query(sqlStatement)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var tag string
		var resourceID int
		err := rows.Scan(&resourceID, &tag)
		if err != nil {
			log.Println(err)
		}
		resources[resourceIndexMap[resourceID]].Tags = append(resources[resourceIndexMap[resourceID]].Tags, tag)
	}
	return resources
}

// GetResourceByID returns a resource with requested ID
func GetResourceByID(id int) Resource {
	resource := Resource{}
	var resourceTagMap map[int][]string
	resourceTagMap = make(map[int][]string)
	resourceTagMap = getResourceTagMap()
	sqlStatement := `
	SELECT * FROM RESOURCE WHERE ID=$1;`
	err := DB.QueryRow(sqlStatement, id).Scan(&resource.ID, &resource.Name, &resource.Description, &resource.Downloads, &resource.Rating, &resource.Github, &resource.Type, &resource.Verified)
	if err != nil {
		return Resource{}
	}
	resource.Tags = resourceTagMap[resource.ID]
	return resource
}

// GetTaskNameFromID returns name from given ID
func GetTaskNameFromID(taskID string) string {
	id, err := strconv.Atoi(taskID)
	if err != nil {
		log.Println(err)
	}
	sqlStatement := `SELECT NAME FROM TASK WHERE ID=$1`
	var taskName string
	err = DB.QueryRow(sqlStatement, id).Scan(&taskName)
	if err != nil {
		return ""
	}
	return taskName
}

// IncrementDownloads will increment the number of downloads
func IncrementDownloads(taskID string) {
	id, err := strconv.Atoi(taskID)
	if err != nil {
		log.Println(err)
	}
	log.Println(id)
	sqlStatement := `UPDATE RESOURCE SET DOWNLOADS = DOWNLOADS + 1 WHERE ID=$1`
	_, err = DB.Exec(sqlStatement, id)
	if err != nil {
		log.Println(err)
	}
}

func updateAverageRating(resourceID int, rating float64) error {
	sqlStatement := `UPDATE RESOURCE SET RATING=$2 WHERE ID=$1`
	_, err := DB.Exec(sqlStatement, resourceID, rating)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetResourceIDFromName will return resource ID from name
func GetResourceIDFromName(name string) (int, error) {
	sqlStatement := `SELECT ID FROM RESOURCE WHERE NAME=$1`
	var resourceID int
	err := DB.QueryRow(sqlStatement, name).Scan(&resourceID)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return resourceID, nil
}

func resourceExists(resourceName string) bool {
	sqlStatement := `SELECT EXISTS(SELECT 1 FROM RESOURCE WHERE NAME=$1 AND VERIFIED=$2)`
	var exists bool
	err := DB.QueryRow(sqlStatement, resourceName, true).Scan(&exists)
	if err != nil {
		log.Println(err)
	}
	return exists
}

// DeleteResource will delete a resource
func DeleteResource(resourceID int) error {
	sqlStatement := `DELETE FROM RESOURCE WHERE ID=$1`
	_, err := DB.Exec(sqlStatement, resourceID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
