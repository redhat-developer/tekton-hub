/*
This file provides utilify funcitons for task APIs
*/
package api

import "github.com/Pipelines-Marketplace/backend/pkg/models"

// query for all tasks
func allTasks() []models.Task {
	tasks := []models.Task{}
	models.DB.Find(&tasks)
	return tasks
}

// query for task with given id
func getTask(param string) models.Task {
	task := models.Task{}
	if param == "" {
		return models.Task{}
	}
	models.DB.Where("name=?", param).First(&task)
	return task
}
