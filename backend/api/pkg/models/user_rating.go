package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

// UserRating represents relationship between User and Rating
type UserRating struct {
	UserID     int `gorm:"primary_key;" json:"user_id"`
	ResourceID int `gorm:"primary_key;" json:"resource_id"`
	Stars      int `json:"stars"`
}

// AddRating add's rating provided by user
func AddRating(db *gorm.DB, userID int, resourceID int, stars int, prevStars int) interface{} {
	sqlStatement := `INSERT INTO USER_RATING(USER_ID,RESOURCE_ID,STARS) VALUES($1,$2,$3)`
	_, err := db.DB().Exec(sqlStatement, userID, resourceID, stars)
	if err != nil {
		log.Println(err)
		return map[string]interface{}{"status": false, "message": "Use PUT method to update existing rating"}
	}
	err = addStars(db, resourceID, stars, prevStars)
	if err != nil {
		return map[string]interface{}{"status": false, "message": "Not able to add stars to ratings table"}
	}
	averageRating := calculateAverageRating(db, resourceID)
	err = updateAverageRating(db, resourceID, averageRating)
	if err != nil {
		return map[string]interface{}{"status": false, "message": "Unable to update average rating in task"}
	}
	return updatedRatings(db, userID, resourceID)
}

// UpdateRating will update existing rating
func UpdateRating(db *gorm.DB, userID int, resourceID int, stars int, prevStars int) UpdatedRatingResponse {
	sqlStatement := `UPDATE USER_RATING SET STARS=$3 WHERE RESOURCE_ID=$2 AND USER_ID=$1`
	_, err := db.DB().Exec(sqlStatement, userID, resourceID, stars)
	if err != nil {
		log.Println(err)
	}
	updateStars(db, resourceID, stars, prevStars)
	averageRating := calculateAverageRating(db, resourceID)
	updateAverageRating(db, resourceID, averageRating)
	return updatedRatings(db, userID, resourceID)
}

// GetUserRating queries for user rating by id
func GetUserRating(db *gorm.DB, userID int, resourceID int) UserRating {
	userRating := UserRating{}
	sqlStatement := `SELECT * FROM USER_RATING WHERE RESOURCE_ID=$2 AND USER_ID=$1`
	err := db.DB().QueryRow(sqlStatement, userID, resourceID).Scan(&userRating.UserID, &userRating.ResourceID, &userRating.Stars)
	if err != nil {
		log.Println(err)
	}
	return userRating
}

func updatedRatings(db *gorm.DB, userID int, resourceID int) UpdatedRatingResponse {
	updatedRatingResponse := UpdatedRatingResponse{}
	rating := Rating{}
	sqlStatement := `SELECT * FROM RATING WHERE RESOURCE_ID=$1`
	err := db.DB().QueryRow(sqlStatement, resourceID).Scan(&rating.ID, &rating.ResourceID, &rating.OneStar, &rating.TwoStar, &rating.ThreeStar, &rating.FourStar, &rating.FiveStar)
	if err != nil {
		log.Println(err)
	}
	updatedRatingResponse.OneStar = rating.OneStar
	updatedRatingResponse.TwoStar = rating.TwoStar
	updatedRatingResponse.ThreeStar = rating.ThreeStar
	updatedRatingResponse.FourStar = rating.FourStar
	updatedRatingResponse.FiveStar = rating.FiveStar
	updatedRatingResponse.ResourceID = resourceID
	updatedRatingResponse.Average = getRatingFromTask(db, resourceID)
	return updatedRatingResponse
}

func getRatingFromTask(db *gorm.DB, resourceID int) float64 {
	var rating float64
	sqlStatement := `SELECT RATING FROM RESOURCE WHERE ID=$1`
	err := db.DB().QueryRow(sqlStatement, resourceID).Scan(&rating)
	if err != nil {
		log.Println(err)
	}
	return rating
}
