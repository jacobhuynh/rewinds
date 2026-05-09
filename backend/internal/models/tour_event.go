package models

import "time"

type TourAlert struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ArtistID  string    `json:"artist_id"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

type TourEvent struct {
	ID         string    `json:"id"`
	ArtistID   string    `json:"artist_id"`
	Source     string    `json:"source"`
	ExternalID string    `json:"external_id"`
	VenueName  string    `json:"venue_name"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	EventDate  time.Time `json:"event_date"`
	TicketURL  *string   `json:"ticket_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
