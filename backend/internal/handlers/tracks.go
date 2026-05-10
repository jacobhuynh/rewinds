package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TracksHandler struct {
	pool *pgxpool.Pool
}

func NewTracksHandler(pool *pgxpool.Pool) *TracksHandler {
	return &TracksHandler{pool: pool}
}

func (h *TracksHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
