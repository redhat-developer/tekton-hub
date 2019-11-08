package api

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/Pipelines-Marketplace/backend/pkg/compress"
	"github.com/Pipelines-Marketplace/backend/pkg/models"
)

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

// GetCompressedFiles returns the created zip file of requestedTask
func GetCompressedFiles(requestedTask string) *os.File {
	dir := "catalog" + "/" + requestedTask
	requestedFiles, err := ioutil.ReadDir("catalog" + "/" + requestedTask + "/")
	if err != nil {
		log.Fatal(err)
	}
	finalZipFile, err := compress.ZipFiles("finalZipFile.zip", requestedFiles, dir)
	if err != nil {
		log.Println(err)
	}
	return finalZipFile
}
