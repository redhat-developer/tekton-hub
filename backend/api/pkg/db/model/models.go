package model

import (
	"github.com/jinzhu/gorm"
)

type (
	// User represents User model in database
	User struct {
		gorm.Model
		Name      string `gorm:"not null;unique"`
		FirstName string
		LastName  string
		Email     string
		Token     string
	}

	Category struct {
		gorm.Model
		Name string `gorm:"size:100;not null;unique"`
		Tags []Tag
	}

	Tag struct {
		gorm.Model
		Name       string `gorm:"not null;unique"`
		Category   Category
		CategoryID int
	}

	Repository struct {
		gorm.Model
		Name       string
		URL        string
		Owner      string
		ContextDir string
	}

	Resource struct {
		gorm.Model
		Name         string
		Type         string
		Downloads    int
		Rating       float64
		Repository   Repository `gorm:"foreignkey:RepositoryID"`
		RepositoryID int
	}

	ResourceVersion struct {
		gorm.Model
		Resource   Resource
		ResourceID int

		Description string
		Version     string
		URL         string
	}

	ResourceTags struct {
		gorm.Model
		Tag        Tag
		TagID      int
		Resource   Resource
		ResourceID int
	}

	ResourceUserRating struct {
		gorm.Model
		UserID     int
		User       User
		Resource   Resource
		ResourceID int
		Rating     int
	}
)
