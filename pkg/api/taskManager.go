/*
This file provides utilify funcitons for task APIs
*/
package api

import (
	"github.com/backend/pkg/models"
)

func allTasks() []Tasks {
	tasks := []Tasks{}
	models.DB.Find(&tasks)
	return tasks
}
