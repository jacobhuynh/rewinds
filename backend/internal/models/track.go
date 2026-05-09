package models

import "time"

type Track struct {
	ID         string    `json:"id"`
	SpotifyID  string    `json:"spotify_id"`
	AlbumID    string    `json:"album_id"`
	ArtistID   string    `json:"artist_id"`
	Name       string    `json:"name"`
	DurationMs int       `json:"duration_ms"`
	CreatedAt  time.Time `json:"created_at"`
}
