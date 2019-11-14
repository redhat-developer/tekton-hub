package models

import (
	"log"
	"strconv"
)

// Task is a database model representing task data
type Task struct {
	ID          int
	Name        string
	Description string
	Downloads   int
	Rating      float64
	Github      string
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
	return tasks
}

// GetTaskWithName returns a task with requested ID
func GetTaskWithName(name string) Task {
	task := Task{}
	id, err := strconv.Atoi(name)
	if err != nil {
		log.Fatalln(err)
	}
	sqlStatement := `
	SELECT * FROM TASK WHERE ID=$1;`
	err = DB.QueryRow(sqlStatement, id).Scan(&task.ID, &task.Name, &task.Description, &task.Downloads, &task.Rating, &task.Github)
	if err != nil {
		return Task{}
	}
	return task
}
