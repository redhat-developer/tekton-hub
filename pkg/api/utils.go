package api

// AddRatingsRequest represents request body for adding ratings
type AddRatingsRequest struct {
	UserID     int `json:"user_id"`
	ResourceID int `json:"resource_id"`
	Stars      int `json:"stars"`
	PrevStars  int `json:"prev_stars"`
}

// OAuthAccessResponse represents access_token
type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

// Code will
type Code struct {
	Token string `json:"token"`
}
