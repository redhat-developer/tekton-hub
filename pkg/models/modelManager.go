package models

import (
	"log"
	"os"
	"strings"

	"github.com/Pipelines-Marketplace/backend/pkg/polling"
	"github.com/Pipelines-Marketplace/backend/pkg/utility"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"

	// Blank for package side effect( Calling init() )
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB is a PostgreSQL gorm object
var DB *gorm.DB

// StartConnection will start a new database connection
func StartConnection() error {
	// Connect to PostgreSQL on Openshift
	db, err := gorm.Open("postgres", "host=localhost port=15432 user=postgres dbname=marketplace password=postgres sslmode=disable")
	// db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=marketplace password=postgres sslmode=disable")
	DB = db
	if err != nil {
		return err
	}
	// defer db.Close()
	err = db.DB().Ping()
	if err != nil {
		return err
	}
	println("Successful Connection")
	return nil
}

// CreateTaskTable creates a new table of type Task{}
func CreateTaskTable() {
	// Migrate to Schema
	DB.CreateTable(&Task{})
}

func storeContentsInFile(dir *github.RepositoryContent, file *github.RepositoryContent, content *string) {
	f, err := os.OpenFile("catalog/"+dir.GetName()+"/"+file.GetName(), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println("Cannot open new file")
	}
	defer f.Close()
	if _, err = f.WriteString(*content); err != nil {
		log.Println("Cannot append to file")
	}
}

func extractREADMEFile(file *github.RepositoryContent, dir *github.RepositoryContent, task *Task) {
	if strings.HasSuffix(file.GetName(), ".md") {
		// Get the contents of README file
		taskDescription, err := utility.GetREADMEContent(dir, file)
		if err != nil {
			log.Fatalln(err)
		}
		storeContentsInFile(dir, file, &taskDescription)
		task.Description = file.GetDownloadURL()
	}
}

func extractYAMLFile(file *github.RepositoryContent, dir *github.RepositoryContent, task *Task) {
	if strings.HasSuffix(file.GetName(), ".yaml") {
		yamlContent, err := utility.GetYAMLContent(dir, file)
		if err != nil {
			log.Fatalln(err)
		}
		storeContentsInFile(dir, file, &yamlContent)
		task.YAML = file.GetDownloadURL()
	}
}

// AddContentsToDB will add contents from Github catalog to database
func AddContentsToDB() {
	task := Task{}
	// Get all directories
	repoContents, err := polling.GetDirContents(utility.Ctx, utility.Client, "tektoncd", "catalog", "", nil)
	if err != nil {
		log.Fatalln(err)
	}
	os.Mkdir("catalog", 0777)
	for _, dir := range repoContents {
		if utility.IsValidDirectory(dir) {
			d, err := polling.GetDirContents(utility.Ctx, utility.Client, "tektoncd", "catalog", dir.GetName(), nil)
			if err != nil {
				log.Fatalln(err)
			}
			task.Name = dir.GetName()
			task.Rating = 4.5
			task.Downloads = 10
			os.Mkdir("catalog/"+dir.GetName(), 0777)
			// Iterate over all files in directory
			for _, file := range d {
				extractREADMEFile(file, dir, &task)
				extractYAMLFile(file, dir, &task)
			}
			println(task.Name)
			DB.Create(&task)
		}
	}
}
