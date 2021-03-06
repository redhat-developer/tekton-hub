package models

import (
	"fmt"
	"log"
)

// Rating represents Rating model in database
type Rating struct {
	ID         int `gorm:"primary_key;" json:"id"`
	ResourceID int `gorm:"primary_key;unique" json:"resource_id"`
	OneStar    int `json:"one" gorm:"default:0"`
	TwoStar    int `json:"two" gorm:"default:0"`
	ThreeStar  int `json:"three" gorm:"default:0"`
	FourStar   int `json:"four" gorm:"default:0"`
	FiveStar   int `json:"five" gorm:"default:0"`
}

// PrevStarRequest represents previous stars
type PrevStarRequest struct {
	UserID     int `json:"user_id"`
	ResourceID int `json:"resource_id"`
}

// GetRatingDetialsByResourceID retrieves rating details of a task
func GetRatingDetialsByResourceID(resourceID int) Rating {
	sqlStatement := `SELECT * FROM RATING WHERE RESOURCE_ID=$1`
	taskRating := Rating{}
	err := DB.QueryRow(sqlStatement, resourceID).Scan(&taskRating.ID, &taskRating.ResourceID, &taskRating.OneStar, &taskRating.TwoStar, &taskRating.ThreeStar, &taskRating.FourStar, &taskRating.FiveStar)
	if err != nil {
		log.Println(err)
	}
	return taskRating
}

func calculateAverageRating(resourceID int) float64 {
	rating := Rating{}
	sqlStatement := `SELECT * FROM RATING WHERE RESOURCE_ID=$1`
	err := DB.QueryRow(sqlStatement, resourceID).Scan(&rating.ID, &rating.ResourceID, &rating.OneStar, &rating.TwoStar, &rating.ThreeStar, &rating.FourStar, &rating.FiveStar)
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
		return "ONE_STAR"
	case 2:
		return "TWO_STAR"
	case 3:
		return "THREE_STAR"
	case 4:
		return "FOUR_STAR"
	case 5:
		return "FIVE_STAR"
	}
	return ""
}

func addStars(taskID int, stars int, prevStars int) error {
	starsString := getStarsInString(stars)
	sqlStatement := fmt.Sprintf("INSERT INTO RATING(%v,RESOURCE_ID) VALUES($1,$2) ON CONFLICT (RESOURCE_ID) DO UPDATE SET %v=RATING.%v+1", starsString, starsString, starsString)
	_, err := DB.Exec(sqlStatement, 1, taskID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func updateStars(resourceID int, stars int, prevStars int) {
	starsString := getStarsInString(stars)
	sqlStatement := fmt.Sprintf("UPDATE RATING SET %v=%v+1 WHERE RESOURCE_ID=$1", starsString, starsString)
	_, err := DB.Exec(sqlStatement, resourceID)
	if err != nil {
		log.Println(err)
	}
	deleteOldStars(resourceID, prevStars)
}

func deleteOldStars(resourceID int, prevStars int) {
	starsString := getStarsInString(prevStars)
	sqlStatement := fmt.Sprintf("UPDATE RATING SET %v=%v-1 WHERE RESOURCE_ID=$1", starsString, starsString)
	_, err := DB.Exec(sqlStatement, resourceID)
	if err != nil {
		log.Println(err)
	}
}
