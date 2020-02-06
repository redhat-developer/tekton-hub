package models

import "log"

// UserRating represents relationship between User and Rating
type UserRating struct {
	UserID     int `json:"user_id"`
	ResourceID int `json:"resource_id"`
	Stars      int `json:"stars"`
}

// AddRating add's rating provided by user
func AddRating(userID int, resourceID int, stars int, prevStars int) interface{} {
	sqlStatement := `INSERT INTO USER_RATING(USER_ID,RESOURCE_ID,STARS) VALUES($1,$2,$3)`
	_, err := DB.Exec(sqlStatement, userID, resourceID, stars)
	if err != nil {
		log.Println(err)
		return map[string]interface{}{"status": false, "message": "Use PUT method to update existing rating"}
	}
	err = addStars(resourceID, stars, prevStars)
	if err != nil {
		return map[string]interface{}{"status": false, "message": "Not able to add stars to ratings table"}
	}
	averageRating := calculateAverageRating(resourceID)
	err = updateAverageRating(resourceID, averageRating)
	if err != nil {
		return map[string]interface{}{"status": false, "message": "Unable to update average rating in task"}
	}
	return updatedRatings(userID, resourceID)
}

// UpdateRating will update existing rating
func UpdateRating(userID int, resourceID int, stars int, prevStars int) UpdatedRatingResponse {
	sqlStatement := `UPDATE USER_RATING SET STARS=$3 WHERE RESOURCE_ID=$2 AND USER_ID=$1`
	_, err := DB.Exec(sqlStatement, userID, resourceID, stars)
	if err != nil {
		log.Println(err)
	}
	updateStars(resourceID, stars, prevStars)
	averageRating := calculateAverageRating(resourceID)
	updateAverageRating(resourceID, averageRating)
	return updatedRatings(userID, resourceID)
}

// GetUserRating queries for user rating by id
func GetUserRating(userID int, resourceID int) UserRating {
	userRating := UserRating{}
	sqlStatement := `SELECT * FROM USER_RATING WHERE RESOURCE_ID=$2 AND USER_ID=$1`
	err := DB.QueryRow(sqlStatement, userID, resourceID).Scan(&userRating.UserID, &userRating.ResourceID, &userRating.Stars)
	if err != nil {
		log.Println(err)
	}
	return userRating
}

func updatedRatings(userID int, resourceID int) UpdatedRatingResponse {
	updatedRatingResponse := UpdatedRatingResponse{}
	rating := Rating{}
	sqlStatement := `SELECT * FROM RATING WHERE RESOURCE_ID=$1`
	err := DB.QueryRow(sqlStatement, resourceID).Scan(&rating.ID, &rating.ResourceID, &rating.OneStar, &rating.TwoStar, &rating.ThreeStar, &rating.FourStar, &rating.FiveStar)
	if err != nil {
		log.Println(err)
	}
	updatedRatingResponse.OneStar = rating.OneStar
	updatedRatingResponse.TwoStar = rating.TwoStar
	updatedRatingResponse.ThreeStar = rating.ThreeStar
	updatedRatingResponse.FourStar = rating.FourStar
	updatedRatingResponse.FiveStar = rating.FiveStar
	updatedRatingResponse.ResourceID = resourceID
	updatedRatingResponse.Average = getRatingFromTask(resourceID)
	return updatedRatingResponse
}

func getRatingFromTask(resourceID int) float64 {
	var rating float64
	sqlStatement := `SELECT RATING FROM RESOURCE WHERE ID=$1`
	err := DB.QueryRow(sqlStatement, resourceID).Scan(&rating)
	if err != nil {
		log.Println(err)
	}
	return rating
}
