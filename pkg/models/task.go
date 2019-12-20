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
}

func addTask(task *Resource) {
	sqlStatement := `
	INSERT INTO TASK (NAME,DESCRIPTION,DOWNLOADS,RATING,GITHUB)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := DB.Exec(sqlStatement, task.Name, task.Description, task.Downloads, task.Rating, task.Github)
	if err != nil {
		log.Println(err)
	}
}

// AddTask will add a new task
func AddTask(task *Resource, userID int) (int, error) {
	var taskID int
	sqlStatement := `
	INSERT INTO TASK (NAME,DESCRIPTION,DOWNLOADS,RATING,GITHUB)
	VALUES ($1, $2, $3, $4, $5) RETURNING ID`
	err := DB.QueryRow(sqlStatement, task.Name, task.Description, task.Downloads, task.Rating, task.Github).Scan(&taskID)
	if err != nil {
		return 0, err
	}
	// Add Tags separately
	if len(task.Tags) > 0 {
		for _, tag := range task.Tags {
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
				addTaskTag(taskID, tagID)
			} else {
				var newTagID int
				newTagID, err = AddTag(tag)
				if err != nil {
					log.Println(err)
				}
				addTaskTag(taskID, newTagID)
			}
		}
	}
	return taskID, addUserTask(userID, taskID)
}

func addTaskTag(taskID int, tagID int) {
	sqlStatement := `INSERT INTO RESOURCE_TAG(TASK_ID,TAG_ID) VALUES($1,$2)`
	_, err := DB.Exec(sqlStatement, taskID, tagID)
	if err != nil {
		log.Println(err)
	}
}

func addUserTask(userID int, taskID int) error {
	sqlStatement := `INSERT INTO USER_TASK(TASK_ID,USER_ID) VALUES($1,$2)`
	_, err := DB.Exec(sqlStatement, taskID, userID)
	if err != nil {
		return err
	}
	return nil
}

// CheckSameTaskUpload will checkif the user submitted the same task again
func CheckSameTaskUpload(userID int, name string) bool {
	sqlStatement := `SELECT T.NAME FROM TASK T JOIN USER_TASK U ON T.ID=U.TASK_ID WHERE U.USER_ID=$1`
	rows, err := DB.Query(sqlStatement, userID)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var taskName string
		err := rows.Scan(&taskName)
		if err != nil {
			log.Println(err)
		}
		if taskName == name {
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
		err = rows.Scan(&resource.ID, &resource.Name, &resource.Description, &resource.Downloads, &resource.Rating, &resource.Github, &resource.Type)
		if err != nil {
			log.Println(err)
		}
		resources = append(resources, resource)
	}
	resourceIndexMap := make(map[int]int)
	sqlStatement = `SELECT ID FROM RESOURCE`
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

	sqlStatement = `SELECT R.ID,TG.NAME FROM TAG TG JOIN RESOURCE_TAG TT ON TT.TAG_ID=TG.ID JOIN RESOURCE R ON R.ID=TT.TASK_ID`
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
	err := DB.QueryRow(sqlStatement, id).Scan(&resource.ID, &resource.Name, &resource.Description, &resource.Downloads, &resource.Rating, &resource.Github, &resource.Type)
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
	log.Println(taskName)
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
		log.Println(name)
		return 0, err
	}
	return resourceID, nil
}
