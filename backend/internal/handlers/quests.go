package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestsHandler struct {
	pool *pgxpool.Pool
}

func NewQuestsHandler(pool *pgxpool.Pool) *QuestsHandler {
	return &QuestsHandler{pool: pool}
}

func (h *QuestsHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
