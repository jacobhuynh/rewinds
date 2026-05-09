package models

import "time"

type Artist struct {
	ID        string    `json:"id"`
	SpotifyID string    `json:"spotify_id"`
	Name      string    `json:"name"`
	ImageURL  *string   `json:"image_url,omitempty"`
	Genres    []string  `json:"genres"`
	Followers int       `json:"followers"`
	Popularity int      `json:"popularity"`
	CreatedAt time.Time `json:"created_at"`
}
