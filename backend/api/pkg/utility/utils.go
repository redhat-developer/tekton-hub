package utility

import (
	"errors"
	"log"
	"strings"

	"github.com/Pipelines-Marketplace/backend/pkg/polling"
	"github.com/google/go-github/github"
)

var client, ctx = polling.Authenticate()

// Client : Github Client
var Client = client

// Ctx : Context
var Ctx = ctx

// IsValidDirectory checks if the directory is a valid catalog directory
func IsValidDirectory(dir *github.RepositoryContent) bool {
	if dir.GetType() == "dir" && dir.GetName() != "vendor" && dir.GetName() != "test" && dir.GetName() != ".github" {
		return true
	}
	return false
}

// GetREADMEContent returns the content of README file
func GetREADMEContent(dir *github.RepositoryContent, file *github.RepositoryContent) (string, error) {
	if strings.HasSuffix(file.GetName(), ".md") {
		desc, err := polling.GetFileContent(Ctx, Client, "tektoncd", "catalog", dir.GetName()+"/"+file.GetName(), nil)
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

// GetYAMLContent returns content of a YAML file
func GetYAMLContent(dir *github.RepositoryContent, file *github.RepositoryContent) (string, error) {
	if strings.HasSuffix(file.GetName(), ".yaml") {
		desc, err := polling.GetFileContent(Ctx, Client, "tektoncd", "catalog", dir.GetName()+"/"+file.GetName(), nil)
		if err != nil {
			log.Fatalln(err)
		}
		content, err := desc.GetContent()
		if err != nil {
			log.Fatalln(err)
		}

		return content, err
	}
	return "", errors.New("Cannot open YAML")
}
