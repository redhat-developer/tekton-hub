package models

import (
	"fmt"
	"log"
	"strconv"
)

// Rating represents Rating model in database
type Rating struct {
	ID        int `json:"id"`
	TaskID    int `json:"task_id"`
	OneStar   int `json:"one"`
	TwoStar   int `json:"two"`
	ThreeStar int `json:"three"`
	FourStar  int `json:"four"`
	FiveStar  int `json:"five"`
}

// GetRatingDetialsByTaskID retrieves rating details of a task
func GetRatingDetialsByTaskID(id string) Rating {
	taskID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
	}
	sqlStatement := `SELECT * FROM RATING WHERE TASK_ID=$1`
	taskRating := Rating{}
	err = DB.QueryRow(sqlStatement, taskID).Scan(&taskRating.ID, &taskRating.TaskID, &taskRating.OneStar, &taskRating.TwoStar, &taskRating.ThreeStar, &taskRating.FourStar, &taskRating.FiveStar)
	if err != nil {
		log.Println(err)
	}
	return taskRating
}

func calculateAverageRating(taskID int) float64 {
	rating := Rating{}
	sqlStatement := `SELECT * FROM RATING WHERE TASK_ID=$1`
	err := DB.QueryRow(sqlStatement, taskID).Scan(&rating.ID, &rating.TaskID, &rating.OneStar, &rating.TwoStar, &rating.ThreeStar, &rating.FourStar, &rating.FiveStar)
	if err != nil {
		log.Println(err)
	}
	totalResponses := float64(rating.OneStar + rating.TwoStar + rating.ThreeStar + rating.FourStar + rating.FiveStar)
	averageRating := float64(rating.OneStar+rating.TwoStar*2+rating.ThreeStar*3+rating.FourStar*4+rating.FiveStar*5) / (totalResponses)
	return averageRating
}

func getStarsInString(stars int) string {
	switch stars {
	case 1:
		return "ONESTAR"
	case 2:
		return "TWOSTAR"
	case 3:
		return "THREESTAR"
	case 4:
		return "FOURSTAR"
	case 5:
		return "FIVESTAR"
	}
	return ""
}

func addStars(taskID int, stars int, prevStars int) error {
	starsString := getStarsInString(stars)
	sqlStatement := fmt.Sprintf("INSERT INTO RATING(%v,TASK_ID) VALUES($1,$2) ON CONFLICT (TASK_ID) DO UPDATE SET %v=RATING.%v+1", starsString, starsString, starsString)
	_, err := DB.Exec(sqlStatement, 1, taskID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func updateStars(taskID int, stars int, prevStars int) {
	starsString := getStarsInString(stars)
	sqlStatement := fmt.Sprintf("UPDATE RATING SET %v=%v+1 WHERE TASK_ID=$1", starsString, starsString)
	_, err := DB.Exec(sqlStatement, taskID)
	if err != nil {
		log.Println(err)
	}
	deleteOldStars(taskID, prevStars)
}

func deleteOldStars(taskID int, prevStars int) {
	starsString := getStarsInString(prevStars)
	sqlStatement := fmt.Sprintf("UPDATE RATING SET %v=%v-1 WHERE TASK_ID=$1", starsString, starsString)
	_, err := DB.Exec(sqlStatement, taskID)
	if err != nil {
		log.Println(err)
	}
}
