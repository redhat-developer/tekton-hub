package service

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/db/model"
	"go.uber.org/zap"
)

type Resource struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

// ResourceDetail abstracts necessary fields for UI
type ResourceDetail struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Catalog       Catalog   `json:"catalog"`
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	Versions      []Version `json:"versions"`
	Tags          []Tag     `json:"tags"`
	Rating        float64   `json:"rating"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

// ResourceVersionDetail abstracts necessary fields for UI
type ResourceVersionDetail struct {
	ResourceID  uint   `json:"resource_id"`
	VersionID   uint   `json:"version_id"`
	Version     string `json:"version"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type ResourceVersion struct {
	ResourceID int
	VersionID  int
}

type Catalog struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
}

type Version struct {
	ID      uint   `json:"id"`
	Version string `json:"version"`
}

type Tag struct {
	ID  uint   `json:"id"`
	Tag string `json:"tag"`
}

type Filter struct {
	Limit int
}

// Init Convert Resource object to ResourceDetails
func (d *ResourceDetail) Init(r *model.Resource) {
	d.ID = r.ID
	d.Name = r.Name
	d.Type = r.Type
	d.Rating = r.Rating

	d.Versions = make([]Version, len(r.Versions))
	for i, v := range r.Versions {
		d.Versions[i].ID = v.ID
		d.Versions[i].Version = v.Version
	}
	d.Tags = make([]Tag, len(r.Tags))
	for i, t := range r.Tags {
		d.Tags[i].ID = t.ID
		d.Tags[i].Tag = t.Name
	}

	d.Catalog.ID = r.Catalog.ID
	d.Catalog.Type = r.Catalog.Type

	latestVersion := r.Versions[len(r.Versions)-1]
	d.Description = latestVersion.Description
	d.LastUpdatedAt = latestVersion.UpdatedAt
}

// All Resources
func (r *Resource) All(filter Filter) ([]ResourceDetail, error) {

	var all []*model.Resource
	r.db.Order("rating desc, name").Limit(filter.Limit).
		Preload("Catalog").Preload("Versions").Preload("Tags").
		Find(&all)

	ret := make([]ResourceDetail, len(all))
	for i, r := range all {
		ret[i].Init(r)
	}
	return ret, nil
}

// Init converts ResourceVersion Object to ResourceVersionDetail
func (d *ResourceVersionDetail) Init(r *model.ResourceVersion, rv ResourceVersion) {
	d.ResourceID = uint(rv.ResourceID)
	d.VersionID = r.ID
	d.Version = r.Version
	d.Description = r.Description
	d.URL = r.URL
}

// ByVersionID Get resource by version Id
func (r *Resource) ByVersionID(rv ResourceVersion) (ResourceVersionDetail, error) {

	var resources []*model.ResourceVersion
	r.db.First(&resources, rv.VersionID)

	if len(resources) == 0 {
		return ResourceVersionDetail{}, errors.New("Record not found")
	}
	var versionDetail ResourceVersionDetail
	versionDetail.Init(resources[0], rv)

	return versionDetail, nil
}
