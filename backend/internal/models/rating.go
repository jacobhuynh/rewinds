package models

import "time"

type Rating struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	ArtistID     string    `json:"artist_id"`
	AlbumID      *string   `json:"album_id,omitempty"`
	TrackID      *string   `json:"track_id,omitempty"`
	Score        float64   `json:"score"`
	RatingMethod string    `json:"rating_method"`
	Annotation   *string   `json:"annotation,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
