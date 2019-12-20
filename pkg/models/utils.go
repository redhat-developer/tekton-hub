package models

// UpdatedRatingResponse represents response for ratings API
type UpdatedRatingResponse struct {
	ResourceID int     `json:"resource_id"`
	OneStar    int     `json:"one_star"`
	TwoStar    int     `json:"two_star"`
	ThreeStar  int     `json:"three_star"`
	FourStar   int     `json:"four_star"`
	FiveStar   int     `json:"five_star"`
	Average    float64 `json:"average"`
}
