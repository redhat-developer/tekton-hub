package upload

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Pipelines-Marketplace/backend/pkg/polling"
	"github.com/Pipelines-Marketplace/backend/pkg/utility"
)

// NewUploadObject  represents new task/pipelines
type NewUploadObject struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Github      string   `json:"github"`
	Tags        []string `json:"tags"`
}

// CodeSearchResult represents the result of a code search.
type CodeSearchResult struct {
	Total             *int         `json:"total_count,omitempty"`
	IncompleteResults *bool        `json:"incomplete_results,omitempty"`
	CodeResults       []CodeResult `json:"items,omitempty"`
}

// CodeResult represents a single search result.
type CodeResult struct {
	Name        *string     `json:"name,omitempty"`
	Path        *string     `json:"path,omitempty"`
	SHA         *string     `json:"sha,omitempty"`
	HTMLURL     *string     `json:"html_url,omitempty"`
	Repository  interface{} `json:"repository,omitempty"`
	TextMatches interface{} `json:"text_matches,omitempty"`
}

func getGithubOwner(githubURL string) (string, string) {
	githubURLElements := strings.Split(githubURL, "/")
	owner := githubURLElements[len(githubURLElements)-2]
	repositoryName := githubURLElements[len(githubURLElements)-1]
	return owner, repositoryName
}

// GetAllYAMLFilesFromRepository will fetch all YAML files from a repository
func GetAllYAMLFilesFromRepository(taskName string, githubURL string) {
	owner, repositoryName := getGithubOwner(githubURL)
	// Check if owner and repository are valid
	query := fmt.Sprintf("https://api.github.com/search/code?q=%v+repo:%v/%v+extension:yaml", taskName, owner, repositoryName)
	log.Println(query)
	var result CodeSearchResult
	request, err := http.NewRequest("GET", query, nil)
	resp, err := utility.Client.Do(utility.Ctx, request, &result)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	// Get the latest commit
	latestCommitSHA := getLatestCommit(owner, repositoryName)
	log.Println(latestCommitSHA)

	paths := getPathsFromCodeResult(result.CodeResults)
	log.Println(paths)
	desc, err := polling.GetFileContent(utility.Ctx, utility.Client, owner, repositoryName, paths[0], nil)
	if err != nil {
		log.Fatalln(err)
	}
	content, err := desc.GetContent()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(content)

}

func getPathsFromCodeResult(CodeResults []CodeResult) []string {
	var filePaths []string
	for _, result := range CodeResults {
		filePaths = append(filePaths, *result.Path)
	}
	return filePaths
}

func isValidGithubURL() {

}

func getLatestCommit(owner string, repositoryName string) string {
	commitInfo, _, err := utility.Client.Repositories.ListCommits(utility.Ctx, owner, repositoryName, nil)
	if err != nil {
		log.Println(err)
	}
	return commitInfo[0].GetSHA()
}

func getContent() {

}
