package models

import "time"

type Prediction struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	ArtistID         string    `json:"artist_id"`
	Genre            string    `json:"genre"`
	CurrentRank      int       `json:"current_rank"`
	TargetRank       int       `json:"target_rank"`
	Deadline         time.Time `json:"deadline"`
	Wager            int       `json:"wager"`
	PayoutMultiplier float64   `json:"payout_multiplier"`
	Outcome          string    `json:"outcome"`
	PredictedAt      time.Time `json:"predicted_at"`
	ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
}
