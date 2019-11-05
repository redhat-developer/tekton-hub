/*
This file provides utilify funcitons for task APIs
*/
package api

import (
	"errors"
	"log"
	"strings"

	"github.com/backend/pkg/polling"
	"github.com/google/go-github/github"
)

// Check for a valid directory
func isValidDirectory(dir *github.RepositoryContent) bool {
	if dir.GetType() != "file" && dir.GetName() != ".github" && dir.GetName() != "LICENSE" && dir.GetName() != "OWNERS" && dir.GetName() != "vendor" && dir.GetName() != "test" {
		return true
	}
	return false
}

// Get description from a README file
func getDescription(dir *github.RepositoryContent, file *github.RepositoryContent) (string, error) {
	if strings.HasSuffix(file.GetName(), ".md") {
		desc, err := polling.GetFileContent(client, ctx, "tektoncd", "catalog", dir.GetName()+"/"+file.GetName(), nil)
		if err != nil {
			log.Fatalln(err)
		}
		content, err := desc.GetContent()
		if err != nil {
			log.Fatalln(err)
		}
		return content, err
	}
	return "", errors.New("Cannot open README")
}

// Get the basic details of all the tasks
func allTasks() []Tasks {
	tasks := []Tasks{}
	task := Tasks{}
	// Get all directories
	repoContents, err := polling.GetDirContents(client, ctx, "tektoncd", "catalog", "", nil)
	if err != nil {
		log.Fatalln(err)
	}
	for index, dir := range repoContents {
		if isValidDirectory(dir) {
			d, err := polling.GetDirContents(client, ctx, "tektoncd", "catalog", dir.GetName(), nil)
			if err != nil {
				log.Fatalln(err)
			}
			task.Name = dir.GetName()
			task.Id = index
			task.Rating = 4.5
			task.Downloads = 10
			// Iterate over all files in directory
			for _, file := range d {
				if strings.HasSuffix(file.GetName(), ".md") {
					// Get the contents of README file
					task.Description, err = getDescription(dir, file)
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
			tasks = append(tasks, task)
		}
	}
	return tasks
}
