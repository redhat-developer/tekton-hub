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

// UpdateResourceRating update user's rating of a resource and resource's average rating
func (r *Rating) UpdateResourceRating(rd RatingDetails) (string, error) {

	r.db.Where("user_id = ? AND resource_id = ?", rd.UserID, rd.ResourceID).
		Assign(&model.UserResourceRating{Rating: rd.ResourceRating}).
		FirstOrCreate(&model.UserResourceRating{
			UserID:     rd.UserID,
			ResourceID: rd.ResourceID,
			Rating:     rd.ResourceRating,
		})

	var avg float64
	r.db.Table("user_resource_ratings").Where("resource_id = ?", rd.ResourceID).
		Select("avg(rating)").Row().Scan(&avg)

	r.db.Model(&model.Resource{}).Where("id = ?", rd.ResourceID).
		Update("rating", math.Round(avg*10)/10)

	return "Rating Updated", nil
}
