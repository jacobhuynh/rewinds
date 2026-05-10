package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PromoHandler struct {
	pool *pgxpool.Pool
}

func NewPromoHandler(pool *pgxpool.Pool) *PromoHandler {
	return &PromoHandler{pool: pool}
}

func (h *PromoHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
