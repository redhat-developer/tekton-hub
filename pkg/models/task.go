package models

import (
	"log"
	"strconv"
)

// Task is a database model representing task data
type Task struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Downloads   int      `json:"downloads"`
	Rating      float64  `json:"rating"`
	Github      string   `json:"github"`
	Tags        []string `json:"tags"`
}

func addTask(task *Task) {
	sqlStatement := `
	INSERT INTO TASK (NAME,DESCRIPTION,DOWNLOADS,RATING,GITHUB)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := DB.Exec(sqlStatement, task.Name, task.Description, task.Downloads, task.Rating, task.Github)
	if err != nil {
		log.Println(err)
	}
}

// AddTask will add a new task
func AddTask(task *Task, userID int) (int, error) {
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
			log.Println(tag)
			if tagExists {
				sqlStatement = `INSERT INTO TASK_TAG(TASK_ID,TAG_ID) VALUES($1,$2)`
				_, err = DB.Exec(sqlStatement, taskID, tagID)
				if err != nil {
					log.Println(err)
				}
			} else {
				var newTagID int
				sqlStatement = `INSERT INTO TAG(NAME) VALUES($1) RETURNING ID`
				err = DB.QueryRow(sqlStatement, tag).Scan(&newTagID)
				if err != nil {
					log.Println(err)
				}
				sqlStatement = `INSERT INTO TASK_TAG(TASK_ID,TAG_ID) VALUES($1,$2)`
				_, err = DB.Exec(sqlStatement, taskID, newTagID)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	return taskID, addUserTask(userID, taskID)
}

func addUserTask(userID int, taskID int) error {
	sqlStatement := `INSERT INTO USER_TASK(TASK_ID,USER_ID) VALUES($1,$2)`
	_, err := DB.Exec(sqlStatement, taskID, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllTasks will return all the tasks
func GetAllTasks() []Task {
	tasks := []Task{}
	sqlStatement := `
	SELECT * FROM TASK`
	rows, err := DB.Query(sqlStatement)
	defer rows.Close()
	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Downloads, &task.Rating, &task.Github)
		if err != nil {
			log.Println(err)
		}
		tasks = append(tasks, task)
	}
	taskIndexMap := make(map[int]int)
	sqlStatement = `SELECT ID FROM TASK`
	rows, err = DB.Query(sqlStatement)
	if err != nil {
		log.Println(err)
	}
	taskIndex := 0
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		taskIndexMap[id] = taskIndex
		taskIndex = taskIndex + 1
	}

	sqlStatement = `SELECT T.ID,TG.NAME FROM TAG TG JOIN TASK_TAG TT ON TT.TAG_ID=TG.ID JOIN TASK T ON T.ID=TT.TASK_ID`
	rows, err = DB.Query(sqlStatement)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var tag string
		var taskID int
		err := rows.Scan(&taskID, &tag)
		if err != nil {
			log.Println(err)
		}
		tasks[taskIndexMap[taskID]].Tags = append(tasks[taskIndexMap[taskID]].Tags, tag)
	}
	return tasks
}

// GetTaskWithName returns a task with requested ID
func GetTaskWithName(name string) Task {
	task := Task{}
	id, err := strconv.Atoi(name)
	if err != nil {
		log.Fatalln(err)
	}
	var taskTagMap map[int][]string
	taskTagMap = make(map[int][]string)
	taskTagMap = getTaskTagMap()
	sqlStatement := `
	SELECT * FROM TASK WHERE ID=$1;`
	err = DB.QueryRow(sqlStatement, id).Scan(&task.ID, &task.Name, &task.Description, &task.Downloads, &task.Rating, &task.Github)
	if err != nil {
		return Task{}
	}
	task.Tags = taskTagMap[task.ID]
	return task
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
	sqlStatement := `UPDATE TASK SET DOWNLOADS = DOWNLOADS + 1 WHERE ID=$1`
	_, err = DB.Exec(sqlStatement, id)
	if err != nil {
		log.Println(err)
	}
}

func updateAverageRating(taskID int, rating float64) {
	sqlStatement := `UPDATE TASK SET RATING=$2 WHERE ID=$1`
	_, err := DB.Exec(sqlStatement, taskID, rating)
	if err != nil {
		log.Println(err)
	}
}
