package models

import "time"

type RatingComparison struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ItemAID   string    `json:"item_a_id"`
	ItemBID   string    `json:"item_b_id"`
	ItemAType string    `json:"item_a_type"`
	WinnerID  string    `json:"winner_id"`
	CreatedAt time.Time `json:"created_at"`
}
