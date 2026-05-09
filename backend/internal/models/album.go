package models

import "time"

type Album struct {
	ID          string    `json:"id"`
	SpotifyID   string    `json:"spotify_id"`
	ArtistID    string    `json:"artist_id"`
	Name        string    `json:"name"`
	ImageURL    *string   `json:"image_url,omitempty"`
	ReleaseDate *string   `json:"release_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
