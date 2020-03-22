package models

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"

	// Blank for package side effect
	_ "github.com/lib/pq"
)

// DB is a PostgreSQL object
var DB *sql.DB

//GDB is a Gorm object
var GDB *gorm.DB

// Connect will start a new database connection
func Connect(app app.Config) error {

	log := app.Logger().With("name", "model")
	conn := app.Database().ConnectionString()

	log.Debugf("connecting to db: %s", conn)

	db, err := gorm.Open("postgres", conn)

	if err != nil {
		return err
	}

	// *gorm.DB Object
	GDB = db

	//*sql.DB Object
	DB = db.DB()

	log.Info("Successfully connection")
	return nil
}

func extractDescriptionFromREADME(readmeFile *github.RepositoryContent, dir *github.RepositoryContent) string {
	file, err := os.Open("catalog/" + dir.GetName() + "/" + readmeFile.GetName())
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	isParagraph := false
	description := ""
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			if isParagraph {
				break
			}
			isParagraph = true
			continue
		} else {
			description = description + scanner.Text()
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return description
}

// AddResourcesFromCatalog will add contents from Github catalog to database
//func AddResourcesFromCatalog(owner string, repositoryName string) {
//log.Println("Adding resources from catalog")
//// Get all directories
//repoContents, err := polling.GetDirContents(utility.Ctx, utility.Client, owner, repositoryName, "", nil)
//if err != nil {
//log.Println(err)
//}
//for _, dir := range repoContents {
//if utility.IsValidDirectory(dir) {
//d, err := polling.GetDirContents(utility.Ctx, utility.Client, owner, repositoryName, dir.GetName(), nil)
//if err != nil {
//log.Println(err)
//}
//// Add the resource to DB
//resource := Resource{
//Name:      dir.GetName(),
//Rating:    0.0,
//Downloads: 0.0,
//Github:    "http://github.com/" + owner + "/" + repositoryName,
//Verified:  true,
//}
//var resourceID int
//resourceID, err = AddCatalogResource(&resource)
//if err != nil {
//log.Println(err)
//}
//addGithubDetails(resourceID, owner, repositoryName, "")
//// Iterate over all files in directory
//for _, file := range d {
//resourcePath := dir.GetName() + "/" + file.GetName()
//if strings.HasSuffix(file.GetName(), ".yaml") {
//// Store the path of file
//updateGithubYAMLDetails(resourceID, resourcePath)
//log.Println(dir.GetName() + " " + file.GetName())
//// Store the raw file path
//rawResourcePath := fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/%v/%v", owner, repositoryName, "master", resourcePath)
//AddResourceRawPath(rawResourcePath, resourceID, "Task")
//} else if strings.HasSuffix(file.GetName(), ".md") {
//// Store the path of README file
//log.Println(dir.GetName() + " " + file.GetName())
//updateGithubREADMEDetails(resourceID, resourcePath)
//}
//}
//}
//}
//log.Println("Done!")
//}

// UpdateResourcesFromCatalog will add contents from Github catalog to database
//func UpdateResourcesFromCatalog(owner string, repositoryName string) {
//// Get all directories
//repoContents, err := polling.GetDirContents(utility.Ctx, utility.Client, owner, repositoryName, "", nil)
//if err != nil {
//log.Println(err)
//}
//for _, dir := range repoContents {
//if utility.IsValidDirectory(dir) {
//d, err := polling.GetDirContents(utility.Ctx, utility.Client, owner, repositoryName, dir.GetName(), nil)
//if err != nil {
//log.Println(err)
//}
//// Add the resource to DB
//resource := Resource{
//Name:      dir.GetName(),
//Rating:    0.0,
//Downloads: 0.0,
//Github:    "http://github.com/" + owner + "/" + repositoryName,
//Verified:  true,
//}
//var resourceID int
//// Check if the resource already exists
//if !resourceExists(resource.Name) {
//resourceID, err = AddCatalogResource(&resource)
//if err != nil {
//log.Println(err)
//}
//// Iterate over all files in directory
//for _, file := range d {
//resourcePath := dir.GetName() + "/" + file.GetName()
//addGithubDetails(resourceID, owner, repositoryName, "")
//if strings.HasSuffix(file.GetName(), ".yaml") {
//// Store the path of file
//updateGithubYAMLDetails(resourceID, resourcePath)
//// Store the raw file path
//rawResourcePath := fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/%v/%v", owner, repositoryName, "master", resourcePath)
//AddResourceRawPath(rawResourcePath, resourceID, "Task")
//} else if strings.HasSuffix(file.GetName(), ".md") {
//// Store the path of README file
//updateGithubREADMEDetails(resourceID, resourcePath)
//}
//}
//}
//}
//}
//}
