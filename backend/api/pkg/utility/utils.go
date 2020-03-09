package utility

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/google/go-github/github"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/polling"
	"go.uber.org/zap"
)

//var client, context.Background() = polling.Authenticate()

//// Client : Github Client
//var Client = client

//// context.Background() : Context
//var context.Background() = context.Background()

type GitHub struct {
	app    app.Config
	log    *zap.SugaredLogger
	client *github.Client
}

func New(app app.Config) *GitHub {

	// Authenticate and return a Github client
	gh := app.GitHub().Client
	return &GitHub{
		app:    app,
		log:    app.Logger().With("name", "github"),
		client: gh,
	}
}

// IsValidDirectory checks if the directory is a valid catalog directory
func (gh *GitHub) IsValidDirectory(dir *github.RepositoryContent) bool {
	if dir.GetType() == "dir" && dir.GetName() != "vendor" && dir.GetName() != "test" && dir.GetName() != ".github" {
		return true
	}
	return false
}

// GetREADMEContent returns the content of README file
func (gh *GitHub) GetREADMEContent(dir *github.RepositoryContent, file *github.RepositoryContent) (string, error) {
	if strings.HasSuffix(file.GetName(), ".md") {
		desc, err := polling.GetFileContent(context.Background(), gh.client, "tektoncd", "catalog", dir.GetName()+"/"+file.GetName(), nil)
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
func (gh *GitHub) GetYAMLContent(dir *github.RepositoryContent, file *github.RepositoryContent) (string, error) {
	if strings.HasSuffix(file.GetName(), ".yaml") {
		desc, err := polling.GetFileContent(context.Background(), gh.client, "tektoncd", "catalog", dir.GetName()+"/"+file.GetName(), nil)
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
