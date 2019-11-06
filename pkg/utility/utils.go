package utility

import (
	"errors"
	"log"
	"strings"

	"github.com/backend/pkg/polling"
	"github.com/google/go-github/github"
)

var Client, Ctx = polling.Authenticate()

// Check for a valid directory
func IsValidDirectory(dir *github.RepositoryContent) bool {
	if dir.GetType() != "file" && dir.GetName() != ".github" && dir.GetName() != "LICENSE" && dir.GetName() != "OWNERS" && dir.GetName() != "vendor" && dir.GetName() != "test" {
		return true
	}
	return false
}

// Get description from a README file
func GetDescription(dir *github.RepositoryContent, file *github.RepositoryContent) (string, error) {
	if strings.HasSuffix(file.GetName(), ".md") {
		desc, err := polling.GetFileContent(Client, Ctx, "tektoncd", "catalog", dir.GetName()+"/"+file.GetName(), nil)
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
