package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ToursHandler struct {
	pool *pgxpool.Pool
}

func NewToursHandler(pool *pgxpool.Pool) *ToursHandler {
	return &ToursHandler{pool: pool}
}

func (h *ToursHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
