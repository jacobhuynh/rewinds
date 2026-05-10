package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LeaderboardHandler struct {
	pool *pgxpool.Pool
}

func NewLeaderboardHandler(pool *pgxpool.Pool) *LeaderboardHandler {
	return &LeaderboardHandler{pool: pool}
}

func (h *LeaderboardHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
