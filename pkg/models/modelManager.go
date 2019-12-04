package models

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Pipelines-Marketplace/backend/pkg/polling"
	"github.com/Pipelines-Marketplace/backend/pkg/utility"
	"github.com/google/go-github/github"

	// Blank for package side effect
	_ "github.com/lib/pq"
)

// DB is a PostgreSQL object
var DB *sql.DB

// StartConnection will start a new database connection
func StartConnection() error {
	var (
		host     = "localhost"
		port     = 15432
		user     = os.Getenv("POSTGRESQL_USERNAME")
		password = os.Getenv("POSTGRESQL_PASSWORD")
		dbname   = os.Getenv("POSTGRESQL_DATABASE")
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Connect to PostgreSQL on Openshift
	db, err := sql.Open("postgres", psqlInfo)
	DB = db
	if err != nil {
		return err
	}
	// defer db.Close()
	err = DB.Ping()
	if err != nil {
		return err
	}
	println("Successful Connection")
	return nil
}

func storeContentsInFile(dir *github.RepositoryContent, file *github.RepositoryContent, content *string, taskID int) {
	// f, err := os.OpenFile("catalog/"+dir.GetName()+"/"+file.GetName(), os.O_WRONLY|os.O_CREATE, 0600)
	os.Mkdir("readme", 0777)
	f, err := os.OpenFile("readme/"+strconv.Itoa(taskID)+".md", os.O_WRONLY|os.O_CREATE, 0600)
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
		taskID, _ := GetTaskIDFromName(task.Name)
		sqlStatement := `INSERT INTO TASK_README(TASK_ID,PATH) VALUES($1,$2)`
		path := "readme/" + strconv.Itoa(taskID) + ".md"
		_, err := DB.Exec(sqlStatement, taskID, path)
		if err != nil {
			log.Println(err)
		}
		taskDescription, err := utility.GetREADMEContent(dir, file)
		if err != nil {
			log.Fatalln(err)
		}
		storeContentsInFile(dir, file, &taskDescription, taskID)
		// task.Description = extractDescriptionFromREADME(file, dir)
	}
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

func storeYAMLContentsInFile(dir *github.RepositoryContent, file *github.RepositoryContent, content *string, taskID int) {
	// f, err := os.OpenFile("catalog/"+dir.GetName()+"/"+file.GetName(), os.O_WRONLY|os.O_CREATE, 0600)
	f, err := os.OpenFile("tekton/"+strconv.Itoa(taskID)+".yaml", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println("Cannot open new file")
	}
	defer f.Close()
	if _, err = f.WriteString(*content); err != nil {
		log.Println("Cannot append to file")
	}
}

func extractYAMLFile(file *github.RepositoryContent, dir *github.RepositoryContent, task *Task) {
	if strings.HasSuffix(file.GetName(), ".yaml") {
		taskID, _ := GetTaskIDFromName(task.Name)
		yamlContent, SHA, err := utility.GetYAMLContentWithSHA(dir, file, taskID)
		if err != nil {
			log.Fatalln(err)
		}
		AddNewSHA(taskID, SHA)
		storeYAMLContentsInFile(dir, file, &yamlContent, taskID)
		task.Github = utility.Client.BaseURL.String()
	}
}

// AddContentsToDB will add contents from Github catalog to database
func AddContentsToDB() {
	task := Task{}
	// Get all directories
	repoContents, err := polling.GetDirContents(utility.Ctx, utility.Client, "tektoncd", "catalog", "", nil)
	if err != nil {
		log.Println(err)
	}
	os.Mkdir("tekton", 0777)
	for _, dir := range repoContents {
		if utility.IsValidDirectory(dir) {
			d, err := polling.GetDirContents(utility.Ctx, utility.Client, "tektoncd", "catalog", dir.GetName(), nil)
			if err != nil {
				log.Fatalln(err)
			}
			task.Name = dir.GetName()
			task.Rating = 0.0
			task.Downloads = 0
			// os.Mkdir("catalog/"+dir.GetName(), 0777)
			// Iterate over all files in directory
			log.Println(task.Name)
			for _, file := range d {
				extractREADMEFile(file, dir, &task)
				extractYAMLFile(file, dir, &task)
			}
		}
	}
}
