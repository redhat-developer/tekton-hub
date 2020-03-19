package models

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"gopkg.in/gormigrate.v1"
)

// CreateAndInitialiseTables will create tables
func CreateAndInitialiseTables(db *gorm.DB, log *zap.SugaredLogger) error {

	migration := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// Add Migration Here
		// If writing a migration for a new table then add the same in InitSchema
	})

	migration.InitSchema(func(db *gorm.DB) error {
		if err := db.AutoMigrate(
			// Add all the tables here
			&Tag{},
			&Category{},
			&UserCredential{},
			&Resource{},
			&ResourceTag{},
			&GithubDetail{},
			&Rating{},
			&ResourceRawPath{},
			&UserRating{},
			&UserResource{},
		).Error; err != nil {
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

		log.Info("Schema initialised successfully !!")

		// Add Data to the Tables
		initialiseTables(db)
		log.Info("Data added successfully !!")

		return nil
	})

	if err := migration.Migrate(); err != nil {
		log.Error(err, "could not migrate")
		return err
	}

	log.Info("Migration ran successfully !!")
	return nil
}
