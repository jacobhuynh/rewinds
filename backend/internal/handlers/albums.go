package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AlbumsHandler struct {
	pool *pgxpool.Pool
}

func NewAlbumsHandler(pool *pgxpool.Pool) *AlbumsHandler {
	return &AlbumsHandler{pool: pool}
}

func (h *AlbumsHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
