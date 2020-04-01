package service

import (
	"github.com/jinzhu/gorm"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
	"go.uber.org/zap"
)

type Service interface {
	Resource() *Resource
	Category() *Category
	Rating() *Rating
	User() *User
}

type ServiceImpl struct {
	app app.Base
	log *zap.SugaredLogger
	db  *gorm.DB
}

func New(base app.Base) *ServiceImpl {
	return &ServiceImpl{
		app: base,
		log: base.Logger().With("name", "db"),
		db:  base.DB(),
	}
}

func (s *ServiceImpl) Resource() *Resource {
	return &Resource{s.db, s.log}
}

func (s *ServiceImpl) Category() *Category {
	return &Category{s.db, s.log}
}

func (s *ServiceImpl) Rating() *Rating {
	return &Rating{s.db, s.log}
}

func (s *ServiceImpl) User() *User {
	return &User{s.db, s.log}
}
