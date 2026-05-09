package models

import "time"

type Comment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Body      string    `json:"body"`
	PageType  *string   `json:"page_type,omitempty"`
	PageID    *string   `json:"page_id,omitempty"`
	RatingID  *string   `json:"rating_id,omitempty"`
	ParentID  *string   `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
