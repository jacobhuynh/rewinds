package models

import "time"

type User struct {
	ID                  string     `json:"id"`
	Username            string     `json:"username"`
	PasswordHash        *string    `json:"-"`
	Email               *string    `json:"email,omitempty"`
	AvatarURL           *string    `json:"avatar_url,omitempty"`
	Bio                 *string    `json:"bio,omitempty"`
	IsPremium           bool       `json:"is_premium"`
	Credits             int        `json:"credits"`
	CreditsResetAt      *time.Time `json:"credits_reset_at,omitempty"`
	Points              int        `json:"points"`
	PointsRank          *int       `json:"points_rank,omitempty"`
	SpotifyID           *string    `json:"spotify_id,omitempty"`
	SpotifyAccessToken  *string    `json:"-"`
	SpotifyRefreshToken *string    `json:"-"`
	SpotifyTokenExpiry  *time.Time `json:"-"`
	OnboardingComplete  bool       `json:"onboarding_complete"`
	CreatedAt           time.Time  `json:"created_at"`
}
