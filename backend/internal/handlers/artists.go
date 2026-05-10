package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ArtistsHandler struct {
	pool *pgxpool.Pool
}

func NewArtistsHandler(pool *pgxpool.Pool) *ArtistsHandler {
	return &ArtistsHandler{pool: pool}
}

func (h *ArtistsHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
