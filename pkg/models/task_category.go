package models

import (
	"fmt"
	"log"
)

// TaskCategory represents many-many between Task and Category
type TaskCategory struct {
	TaskID int `json:"taskID"`
	TagID  int `json:"categoryID"`
}

// GetAllTasksWithGivenCategory queries for tags with given categories
func GetAllTasksWithGivenCategory(categories []string) []Task {
	tasks := []Task{}
	args := make([]interface{}, len(categories))
	for index, value := range categories {
		args[index] = value
	}
	params := `$1`
	for index := 1; index <= len(categories); index++ {
		if index > 1 {
			params = params + fmt.Sprintf(",$%d", index)
		}
	}
	log.Println(params)
	var taskCategoryMap map[int][]string
	taskCategoryMap = make(map[int][]string)
	taskCategoryMap = getTaskCategoryMap()
	sqlStatement := `
	SELECT DISTINCT T.ID,T.NAME,T.DESCRIPTION,T.DOWNLOADS,T.RATING,T.GITHUB
	FROM TASK AS T JOIN TASK_CATEGORY AS TC ON (T.ID=TC.TASK_ID) JOIN CATEGORY
	AS CT ON (CT.ID=TC.CATEGORY_ID AND CT.NAME in (` +
		params + `));`
	log.Println(sqlStatement)
	rows, err := DB.Query(sqlStatement, args...)
	defer rows.Close()
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Downloads, &task.Rating, &task.Github)
		if err != nil {
			log.Println(err)
		}
		task.Tags = taskCategoryMap[task.ID]
		tasks = append(tasks, task)
	}
	return tasks
}

func getTaskCategoryMap() map[int][]string {
	sqlStatement := `SELECT DISTINCT T.ID,TG.NAME FROM TASK AS T JOIN 
	TASK_CATEGORY AS TT ON (T.ID=TT.TASK_ID) JOIN CATEGORY AS TG 
	ON (TG.ID=TT.CATEGORY_ID);`
	rows, err := DB.Query(sqlStatement)
	// mapping task ID with tag names
	var taskCategoryMap map[int][]string
	taskCategoryMap = make(map[int][]string)
	for rows.Next() {
		var taskID int
		var categoryName string
		err = rows.Scan(&taskID, &categoryName)
		if err != nil {
			log.Println(err)
		}
		taskCategoryMap[taskID] = append(taskCategoryMap[taskID], categoryName)
	}
	return taskCategoryMap
}
