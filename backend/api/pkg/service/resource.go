package service

import (
	"github.com/jinzhu/gorm"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/db/model"
	"go.uber.org/zap"
)

type Resource struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

type ResourceDetail struct {
	ID       uint   `json:"id"`
	Source   string // official, verified, community
	Type     string
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
	Tags     []string `json:"tags"`
}

func (d *ResourceDetail) Init(r *model.Resource) {
	d.ID = r.ID
	d.Name = r.Name

	d.Versions = make([]string, len(r.Versions))
	for i, v := range r.Versions {
		d.Versions[i] = v.Version
	}
	d.Tags = make([]string, len(r.Tags))
	for i, t := range r.Tags {
		d.Tags[i] = t.Name
	}
}

type Filter struct {
	Limit uint
}

func (r *Resource) All(f Filter) ([]ResourceDetail, error) {
	all := []*model.Resource{}
	// TODO load tags and versions
	// sort by rating and then by name
	//     db.Order("name DESC")
	//     db.Order("name DESC", true) // reorder
	//     db.Order(gorm.Expr("name = ? DESC", "first")) // sql expression
	r.db.Model(&model.Resource{}).Find(all).Limit(f.Limit)

	ret := make([]ResourceDetail, len(all))
	for i, r := range all {
		ret[i].Init(r)
	}
	return ret, nil
}
