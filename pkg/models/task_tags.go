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

// GetAllTasksWithGivenTags queries for all tasks with given tags
func GetAllTasksWithGivenTags(tags []string) []Task {
	tasks := []Task{}
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
	var taskTagMap map[int][]string
	taskTagMap = make(map[int][]string)
	taskTagMap = getTaskTagMap()
	sqlStatement := `
	SELECT DISTINCT T.ID,T.NAME,T.DESCRIPTION,T.DOWNLOADS,T.RATING,T.GITHUB
	FROM TASK AS T JOIN TASK_TAG AS TT ON (T.ID=TT.TASK_ID) JOIN TAG
	AS TG ON (TG.ID=TT.TAG_ID AND TG.NAME in (` +
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
		task.Tags = taskTagMap[task.ID]
		tasks = append(tasks, task)
	}
	return tasks
}

func getTaskTagMap() map[int][]string {
	sqlStatement := `SELECT DISTINCT T.ID,TG.NAME FROM TASK AS T JOIN TASK_TAG AS TT ON (T.ID=TT.TASK_ID) JOIN TAG AS TG ON (TG.ID=TT.TAG_ID);`
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

// AddTagsToTasks will add tags to tasks
func AddTagsToTasks(taskName string, tagName string) {
	var taskID int
	var tagID int
	sqlStatement := `SELECT ID FROM TASK WHERE NAME=$1`
	_ = DB.QueryRow(sqlStatement, taskName).Scan(&taskID)
	sqlStatement = `SELECT ID FROM TAG WHERE NAME=$1`
	_ = DB.QueryRow(sqlStatement, tagName).Scan(&tagID)
	sqlStatement = `INSERT INTO TASK_TAG(TASK_ID,TAG_ID) VALUES($1,$2)`
	_, err := DB.Exec(sqlStatement, taskID, tagID)
	if err != nil {
		log.Println(err)
	}
}
