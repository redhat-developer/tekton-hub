package model

import (
	"github.com/jinzhu/gorm"
)

type (
	Category struct {
		gorm.Model
		Name string `gorm:"size:100;not null;unique"`
		Tags []Tag
	}

	Tag struct {
		gorm.Model
		Name       string `gorm:"size:100;not null;unique"`
		Category   Category
		CategoryID int
	}

	Repository struct {
		gorm.Model
		Name       string
		URL        string
		Owner      string
		ContextDir string
		Resources  []Resource
	}

	Resource struct {
		gorm.Model
		Name             string
		Type             string
		Downloads        uint
		Rating           float64
		RepositoryID     uint
		ResourceVersions []ResourceVersion
		Tags             []Tag `gorm:"many2many:resource_tags;"`
	}

	ResourceVersion struct {
		gorm.Model
		Description string
		Version     string
		URL         string
		Resource    Resource
		ResourceID  uint
	}

	// ResourceTags struct {
	// 	gorm.Model
	// 	Tag        Tag
	// 	TagID      int
	// 	Resource   Resource
	// 	ResourceID int
	// }

	// User represents User model in database
	User struct {
		gorm.Model
		Name      string `gorm:"not null;unique"`
		FirstName string
		LastName  string
		Email     string
		Token     string
	}

	// ResourceUserRating represents User's rating to resources
	ResourceUserRating struct {
		gorm.Model
		UserID     uint
		User       User
		Resource   Resource
		ResourceID uint
		Rating     uint
	}
)
