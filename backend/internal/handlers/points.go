package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PointsHandler struct {
	pool *pgxpool.Pool
}

func NewPointsHandler(pool *pgxpool.Pool) *PointsHandler {
	return &PointsHandler{pool: pool}
}

func (h *PointsHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
