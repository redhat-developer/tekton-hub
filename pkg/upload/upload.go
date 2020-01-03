package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Pipelines-Marketplace/backend/pkg/models"
	"github.com/Pipelines-Marketplace/backend/pkg/polling"
	"github.com/Pipelines-Marketplace/backend/pkg/utility"
	"github.com/ghodss/yaml"
	"github.com/google/go-github/github"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"golang.org/x/oauth2"
)

// NewUploadRequestObject represents new task/pipelines
type NewUploadRequestObject struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Github      string   `json:"github"`
	Tags        []string `json:"tags"`
	UserID      int      `json:"user_id"`
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

// GetGithubOwner will return github owner and repo name from URL
func GetGithubOwner(githubURL string) (string, string) {
	githubURLElements := strings.Split(githubURL, "/")
	owner := githubURLElements[len(githubURLElements)-2]
	repositoryName := githubURLElements[len(githubURLElements)-1]
	return owner, repositoryName
}

func getPathsFromCodeResult(CodeResults []CodeResult) []string {
	var filePaths []string
	for _, result := range CodeResults {
		filePaths = append(filePaths, *result.Path)
	}
	return filePaths
}

func getLatestCommit(owner string, repositoryName string) string {
	commitInfo, _, err := utility.Client.Repositories.ListCommits(utility.Ctx, owner, repositoryName, nil)
	if err != nil {
		log.Println(err)
	}
	return commitInfo[0].GetSHA()
}

// NewUpload handles uploading of new task/pipeline
func NewUpload(name string, description string, objectType string, tags []string, github string, userID int) interface{} {
	isSameResource := models.CheckSameResourceUpload(userID, name)
	if isSameResource {
		return map[string]interface{}{"status": false, "message": objectType + " already exists"}
	}
	// Get owner and repository name from github link
	owner, repositoryName := GetGithubOwner(github)
	// Check if owner and repository name is valid
	paths, err := search(owner, repositoryName, objectType, name, userID)
	if err != nil {
		log.Println(err)
		return map[string]interface{}{"status": false, "message": "The listed users and repositories cannot be searched either because the resources do not exist or you do not have permission to view them."}
	}
	// Check for field name and kind
	var content *string
	// var SHA string
	isTaskPresent := false
	var resourcePath string
	for _, path := range paths {
		content = getObjectContent(path, owner, repositoryName)
		taskJSON, err := yaml.YAMLToJSON([]byte(*content))
		if err != nil {
			log.Println(err)
			return map[string]interface{}{"status": false, "message": "Invalid YAML format"}
		}
		var objmap map[string]interface{}
		if err := json.Unmarshal(taskJSON, &objmap); err != nil {
			log.Println(err)
			return map[string]interface{}{"status": false, "message": "Invalid YAML format"}
		}
		nameMap := objmap["metadata"].(map[string]interface{})
		// Change here for pipeline
		if objmap["kind"] == "Task" && nameMap["name"] == name {
			isTaskPresent = true
			resourcePath = path
			break
		}
	}
	if isTaskPresent == false {
		return map[string]interface{}{"status": false, "message": "Task with the given name doesn't exist"}
	}
	// Perform lint validation and schema validation here

	// Add Task details to DB
	resource := models.Resource{
		Name:        name,
		Github:      github,
		Description: description,
		Tags:        tags,
	}
	rawResourcePath := fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/%v/%v", owner, repositoryName, "master", resourcePath)
	resourceID, err := models.AddResource(&resource, userID, owner, repositoryName, resourcePath)
	if err != nil {
		log.Println(err)
		return map[string]interface{}{"status": false, "message": err}
	}

	// Add a raw path
	models.AddResourceRawPath(rawResourcePath, resourceID, objectType)

	return map[string]interface{}{"status": true, "message": "Upload Successfull"}
}

func doesResourceExist(paths []string, owner string, repositoryName string, resourceName string, objectType string) (bool, string, *string) {
	isResourcePresent := false
	var resourcePath string
	var content *string
	for _, path := range paths {
		content = getObjectContent(path, owner, repositoryName)
		var pipeline v1alpha1.Pipeline
		err := yaml.Unmarshal([]byte(*content), &pipeline)
		if err != nil {
			log.Println("Invalid Resource schema")
			return false, "", nil
		}
		if pipeline.TypeMeta.Kind == objectType && pipeline.ObjectMeta.Name == resourceName {
			isResourcePresent = true
			resourcePath = path
			break
		}
	}
	return isResourcePresent, resourcePath, content
}

// NewUploadPipeline handles uploading of new task/pipeline
func NewUploadPipeline(name string, description string, objectType string, tags []string, github string, userID int) interface{} {
	// isSameResource := models.CheckSameResourceUpload(userID, name)
	// if isSameResource {
	// 	return map[string]interface{}{"status": false, "message": objectType + " already exists"}
	// }
	// Get owner and repository name from github link
	owner, repositoryName := GetGithubOwner(github)
	// Check if owner and repository name is valid
	paths, err := search(owner, repositoryName, objectType, name, userID)
	if err != nil {
		log.Println("Invalid owner and repository name")
		return map[string]interface{}{"status": false, "message": "The listed users and repositories cannot be searched either because the resources do not exist or you do not have permission to view them."}
	}
	// Check for field name and kind
	var content *string
	// var SHA string
	isPipelinePresent := false
	var resourcePath string
	// Check if the resource exists
	isPipelinePresent, resourcePath, content = doesResourceExist(paths, owner, repositoryName, name, objectType)
	if isPipelinePresent == false {
		return map[string]interface{}{"status": false, "message": name + ": Pipeline with the given name doesn't exist"}
	}
	log.Println(resourcePath)
	var pipeline v1alpha1.Pipeline
	err = yaml.Unmarshal([]byte(*content), &pipeline)
	if err != nil {
		fmt.Println("Invalid Pipeline schema")
		return err
	}
	var rawTaskPaths []string
	for _, pipelineTask := range pipeline.Spec.Tasks {
		// For each task get the path of the file
		paths, err := search(owner, repositoryName, "Task", pipelineTask.TaskRef.Name, userID)
		if err != nil {
			fmt.Println("Invalid")
			return map[string]interface{}{"status": false, "message": pipelineTask.TaskRef.Name + ": Task with the given name doesn't exist"}
		}
		isTaskPresent := false
		var taskPath string
		isTaskPresent, taskPath, _ = doesResourceExist(paths, owner, repositoryName, pipelineTask.TaskRef.Name, "Task")
		if isTaskPresent == false {
			return map[string]interface{}{"status": false, "message": name + ": Task with the given name doesn't exist"}
		} else if isTaskPresent == false && taskPath == "" {
			return map[string]interface{}{"status": false, "message": "Invalid Task schema"}
		}
		rawTaskPath := fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/%v/%v", owner, repositoryName, "master", taskPath)
		rawTaskPaths = append(rawTaskPaths, rawTaskPath)
	}
	log.Println(rawTaskPaths)
	// Perform lint validation and schema validation here

	// Add Pipeline details to DB
	resource := models.Resource{
		Name:        name,
		Github:      github,
		Description: description,
		Tags:        tags,
		Type:        objectType,
	}
	rawResourcePath := fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/%v/%v", owner, repositoryName, "master", resourcePath)
	resourceID, err := models.AddResource(&resource, userID, owner, repositoryName, resourcePath)
	if err != nil {
		log.Println(err)
		return map[string]interface{}{"status": false, "message": err}
	}

	// Add a raw path for resource
	models.AddResourceRawPath(rawResourcePath, resourceID, objectType)

	// Add raw paths of pipelines
	for _, rawPath := range rawTaskPaths {
		models.AddResourceRawPath(rawPath, resourceID, "Task")
	}
	return map[string]interface{}{"status": true, "message": "Upload Successfull"}
}

func createTaskFiles(taskID int, name string, content *string) {
	f, err := os.OpenFile("tekton/"+strconv.Itoa(taskID)+".yaml", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err = f.WriteString(*content); err != nil {
		log.Println(err)

	}
}

func getObjectContent(path string, owner string, repositoryName string) *string {
	desc, err := polling.GetFileContent(utility.Ctx, utility.Client, owner, repositoryName, path, nil)
	if err != nil {
		fmt.Println(err)
	}
	content, err := desc.GetContent()
	if err != nil {
		fmt.Println(err)
	}
	return &content
}

// Call search method with a given query
func search(owner string, repositoryName string, objectType string, resourceName string, userID int) ([]string, error) {
	// Use go-github's code search function
	query := fmt.Sprintf("https://api.github.com/search/code?q=kind:%v+%v+repo:%v/%v+extension:yaml", objectType, resourceName, owner, repositoryName)
	var result CodeSearchResult
	request, err := http.NewRequest("GET", query, nil)
	client, ctx := getGithubClientForUser(userID)
	var resp *github.Response
	if client != nil && ctx != nil {
		resp, err = client.Do(ctx, request, &result)
	} else {
		resp, err = utility.Client.Do(utility.Ctx, request, &result)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	paths := getPathsFromCodeResult(result.CodeResults)
	return paths, nil
}

func getGithubClientForUser(userID int) (*github.Client, context.Context) {
	sqlStatement := `SELECT TOKEN FROM USER_CREDENTIAL WHERE ID=$1`
	var token string
	err := models.DB.QueryRow(sqlStatement, userID).Scan(&token)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client, ctx
}
