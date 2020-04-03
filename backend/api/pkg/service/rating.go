package service

import (
	"math"

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

type UpdateRatingDetails struct {
	UserID         uint `json:"user_id"`
	ResourceID     uint `json:"resource_id"`
	ResourceRating uint `json:"rating"`
}

type RatingDetails struct {
	ResourceRating uint `json:"rating"`
}

// GetResourceRating returns user's rating of a resource
func (r *Rating) GetResourceRating(ur UserResource) (RatingDetails, error) {

	rating := &model.UserResourceRating{}
	if r.db.Where("user_id = ? AND resource_id = ?", ur.UserID, ur.ResourceID).Find(&rating).RecordNotFound() {
		return RatingDetails{
			ResourceRating: 0,
		}, nil
	}
	var resRating RatingDetails
	resRating.ResourceRating = rating.Rating

	return resRating, nil
}

// UpdateResourceRating update user's rating of a resource and resource's average rating
func (r *Rating) UpdateResourceRating(rd UpdateRatingDetails) {

	r.db.Where("user_id = ? AND resource_id = ?", rd.UserID, rd.ResourceID).
		Assign(&model.UserResourceRating{Rating: rd.ResourceRating}).
		FirstOrCreate(&model.UserResourceRating{
			UserID:     rd.UserID,
			ResourceID: rd.ResourceID,
			Rating:     rd.ResourceRating,
		})

	var avg float64
	r.db.Model(&model.UserResourceRating{}).Where("resource_id = ?", rd.ResourceID).
		Select("avg(rating)").Row().Scan(&avg)

	r.db.Model(&model.Resource{}).Where("id = ?", rd.ResourceID).
		Update("rating", math.Round(avg*10)/10)

	return
}
