package models

import "log"

// UserRating represents relationship between User and Rating
type UserRating struct {
	UserID int `json:"user_id"`
	TaskID int `json:"task_id"`
	Stars  int `json:"stars"`
}

// AddRating add's rating provided by user
func AddRating(userID int, taskID int, stars int, prevStars int) interface{} {
	sqlStatement := `INSERT INTO USER_RATING(USER_ID,TASK_ID,STARS) VALUES($1,$2,$3)`
	_, err := DB.Exec(sqlStatement, userID, taskID, stars)
	if err != nil {
		log.Println(err)
		return map[string]interface{}{"status": false, "message": "Use PUT method to update existing rating"}
	}
	err = addStars(taskID, stars, prevStars)
	if err != nil {
		return map[string]interface{}{"status": false, "message": "Not able to add stars to ratings table"}
	}
	averageRating := calculateAverageRating(taskID)
	err = updateAverageRating(taskID, averageRating)
	if err != nil {
		return map[string]interface{}{"status": false, "message": "Unable to update average rating in task"}
	}
	return updatedRatings(userID, taskID)
}

// UpdateRating will update existing rating
func UpdateRating(userID int, taskID int, stars int, prevStars int) UpdatedRatingResponse {
	sqlStatement := `UPDATE USER_RATING SET STARS=$3 WHERE TASK_ID=$2 AND USER_ID=$1`
	_, err := DB.Exec(sqlStatement, userID, taskID, stars)
	if err != nil {
		log.Println(err)
	}
	updateStars(taskID, stars, prevStars)
	averageRating := calculateAverageRating(taskID)
	updateAverageRating(taskID, averageRating)
	return updatedRatings(userID, taskID)
}

// GetUserRating queries for user rating by id
func GetUserRating(userID int, taskID int) UserRating {
	userRating := UserRating{}
	sqlStatement := `SELECT * FROM USER_RATING WHERE TASK_ID=$2 AND USER_ID=$1`
	err := DB.QueryRow(sqlStatement, userID, taskID).Scan(&userRating.UserID, &userRating.TaskID, &userRating.Stars)
	if err != nil {
		log.Println(err)
	}
	return userRating
}

func updatedRatings(userID int, taskID int) UpdatedRatingResponse {
	updatedRatingResponse := UpdatedRatingResponse{}
	rating := Rating{}
	sqlStatement := `SELECT * FROM RATING WHERE TASK_ID=$1`
	err := DB.QueryRow(sqlStatement, taskID).Scan(&rating.ID, &rating.TaskID, &rating.OneStar, &rating.TwoStar, &rating.ThreeStar, &rating.FourStar, &rating.FiveStar)
	if err != nil {
		log.Println(err)
	}
	updatedRatingResponse.OneStar = rating.OneStar
	updatedRatingResponse.TwoStar = rating.TwoStar
	updatedRatingResponse.ThreeStar = rating.ThreeStar
	updatedRatingResponse.FourStar = rating.FourStar
	updatedRatingResponse.FiveStar = rating.FiveStar
	updatedRatingResponse.TaskID = taskID
	updatedRatingResponse.Average = getRatingFromTask(taskID)
	return updatedRatingResponse
}

func getRatingFromTask(taskID int) float64 {
	var rating float64
	sqlStatement := `SELECT RATING FROM TASK WHERE ID=$1`
	err := DB.QueryRow(sqlStatement, taskID).Scan(&rating)
	if err != nil {
		log.Println(err)
	}
	return rating
}
