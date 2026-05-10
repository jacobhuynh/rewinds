package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PredictionsHandler struct {
	pool *pgxpool.Pool
}

func NewPredictionsHandler(pool *pgxpool.Pool) *PredictionsHandler {
	return &PredictionsHandler{pool: pool}
}

func (h *PredictionsHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
