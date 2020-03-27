package service

import (
	"github.com/jinzhu/gorm"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/db/model"
	"go.uber.org/zap"
)

// Category Service
type Category struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

// CategoryDetail struct to be returned to UI
type CategoryDetail struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
}

// Init converts Category Object to CategoryDetail
func (d *CategoryDetail) Init(c *model.Category) {
	d.ID = c.ID
	d.Name = c.Name
	d.Tags = make([]Tag, len(c.Tags))
	for i, t := range c.Tags {
		d.Tags[i].ID = t.ID
		d.Tags[i].Tag = t.Name
	}
}

// All Categories with their tags
func (c *Category) All() ([]CategoryDetail, error) {

	var all []*model.Category
	c.db.Preload("Tags").Find(&all)

	ret := make([]CategoryDetail, len(all))
	for i, r := range all {
		ret[i].Init(r)
	}
	return ret, nil
}
