package api

// AddRatingsRequest represents request body for adding ratings
type AddRatingsRequest struct {
	UserID    int `json:"user_id"`
	TaskID    int `json:"task_id"`
	Stars     int `json:"stars"`
	PrevStars int `json:"prev_stars"`
}
