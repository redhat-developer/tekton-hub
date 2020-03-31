package models

import (
	"log"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

// CreateAndInitialiseTables will create tables
func CreateAndInitialiseTables(db *gorm.DB) error {

	gormigrateObj := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// Add Migration Here
		// If writing a migration for a new table then add the same in InitSchema
	})

	gormigrateObj.InitSchema(func(db *gorm.DB) error {
		err := db.AutoMigrate(
			// Add all the tables here
			&Resource{},
			&Category{},
			&Tag{},
			&ResourceTag{},
			&GithubDetail{},
			&Rating{},
			&ResourceRawPath{},
			&UserCredential{},
			&UserRating{},
			&UserResource{},
		).Error

		if err != nil {
			return err
		}

		if err := db.Model(GithubDetail{}).AddForeignKey("resource_id", "resource (id)", "CASCADE", "CASCADE").Error; err != nil {
			return err
		}

		if err := db.Model(Rating{}).AddForeignKey("resource_id", "resource (id)", "CASCADE", "CASCADE").Error; err != nil {
			return err
		}

		if err := db.Model(ResourceRawPath{}).AddForeignKey("resource_id", "resource (id)", "CASCADE", "CASCADE").Error; err != nil {
			return err
		}

		if err := db.Model(ResourceTag{}).AddForeignKey("tag_id", "tag (id)", "CASCADE", "CASCADE").Error; err != nil {
			return err
		}

		if err := db.Model(ResourceTag{}).AddForeignKey("resource_id", "resource (id)", "CASCADE", "CASCADE").Error; err != nil {
			return err
		}

		if err := db.Model(UserResource{}).AddForeignKey("resource_id", "resource (id)", "CASCADE", "CASCADE").Error; err != nil {
			return err
		}

		if err := db.Model(UserRating{}).AddForeignKey("resource_id", "resource (id)", "CASCADE", "CASCADE").Error; err != nil {
			return err
		}

		log.Printf("Schema initialised successfully !!")

		// Add Data to the Tables
		initialiseTables(db)

		log.Printf("Data added successfully !!")

		return nil
	})

	if err := gormigrateObj.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	} else {
		log.Printf("Migration did run successfully !!")
	}

	return nil
}
