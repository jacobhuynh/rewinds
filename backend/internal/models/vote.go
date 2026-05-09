package models

import "time"

type Vote struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	WinnerID  string    `json:"winner_id"`
	LoserID   string    `json:"loser_id"`
	Genre     string    `json:"genre"`
	CreatedAt time.Time `json:"created_at"`
}
