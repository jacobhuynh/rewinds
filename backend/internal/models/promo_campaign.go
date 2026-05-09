package models

import "time"

type PromoCampaign struct {
	ID            string    `json:"id"`
	ArtistID      string    `json:"artist_id"`
	LabelName     string    `json:"label_name"`
	BudgetPts     int       `json:"budget_pts"`
	PtsPerRating  int       `json:"pts_per_rating"`
	RatingsTarget int       `json:"ratings_target"`
	RatingsCount  int       `json:"ratings_count"`
	Genre         string    `json:"genre"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	EndsAt        time.Time `json:"ends_at"`
}
