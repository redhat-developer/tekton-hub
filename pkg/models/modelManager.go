package models

import (
	"log"
	"strings"

	"github.com/backend/pkg/polling"
	"github.com/backend/pkg/utility"
	"github.com/jinzhu/gorm"

	// Blank for package side effect( Calling init() )
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func CreateDatabase() {
	// Connect to PostgreSQL on Openshift
	db, err := gorm.Open("postgres", "host=localhost port=15432 user=postgres dbname=marketplace password=postgres sslmode=disable")
	DB = db
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}
	println("Successful Connection")
	// Migrate to Schema
	// db.CreateTable(&Task{})
}

func AddContentsToDB() {
	task := Task{}
	// Get all directories
	repoContents, err := polling.GetDirContents(utility.Client, utility.Ctx, "tektoncd", "catalog", "", nil)
	if err != nil {
		log.Fatalln(err)
	}
	for _, dir := range repoContents {
		if utility.IsValidDirectory(dir) {
			d, err := polling.GetDirContents(utility.Client, utility.Ctx, "tektoncd", "catalog", dir.GetName(), nil)
			if err != nil {
				log.Fatalln(err)
			}
			task.Name = dir.GetName()
			task.Rating = 4.5
			task.Downloads = 10
			// Iterate over all files in directory
			for _, file := range d {
				if strings.HasSuffix(file.GetName(), ".md") {
					// Get the contents of README file
					task.Description, err = utility.GetDescription(dir, file)
					if err != nil {
						log.Fatalln(err)
					}
					task.Description = file.GetDownloadURL()
				}
				if strings.HasSuffix(file.GetName(), ".yaml") {
					task.YAML = file.GetDownloadURL()
				}
			}
			println(task.Name)
			DB.Create(&task)
		}
	}
}
