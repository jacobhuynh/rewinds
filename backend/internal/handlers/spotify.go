package handlers

import (
	"net/http"

	"github.com/jacobhuynh/rewinds/backend/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// SpotifyHandler handles both Spotify OAuth and Spotify data endpoints.
type SpotifyHandler struct {
	pool  *pgxpool.Pool
	cfg   *config.Config
	redis *redis.Client
}

func NewSpotifyHandler(pool *pgxpool.Pool, cfg *config.Config, redis *redis.Client) *SpotifyHandler {
	return &SpotifyHandler{pool: pool, cfg: cfg, redis: redis}
}

// --- Auth ---

func (h *SpotifyHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *SpotifyHandler) Callback(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *SpotifyHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// --- Data ---

func (h *SpotifyHandler) Onboarding(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *SpotifyHandler) Playlists(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *SpotifyHandler) PlaylistTracks(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
