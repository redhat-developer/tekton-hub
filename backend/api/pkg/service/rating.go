package service

import (
	"github.com/jinzhu/gorm"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/db/model"
	"go.uber.org/zap"
)

// Rating Service
type Rating struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

type UserResource struct {
	UserID     int
	ResourceID int
}

type RatingDetails struct {
	UserID         uint `json:"user_id"`
	ResourceID     uint `json:"resource_id"`
	ResourceRating uint `json:"rating"`
}

func (r *RatingDetails) Init(rating *model.UserResourceRating) {
	r.UserID = rating.UserID
	r.ResourceID = rating.ResourceID
	r.ResourceRating = rating.Rating
}

// GetResourceRating returns user's rating of a resource
func (r *Rating) GetResourceRating(ur UserResource) (RatingDetails, error) {

	rating := &model.UserResourceRating{}
	if r.db.Where("user_id = ? AND resource_id = ?", ur.UserID, ur.ResourceID).Find(&rating).RecordNotFound() {
		return RatingDetails{
			UserID:         uint(ur.UserID),
			ResourceID:     uint(ur.ResourceID),
			ResourceRating: 0,
		}, nil
	}
	var resRating RatingDetails
	resRating.Init(rating)

	return resRating, nil
}
