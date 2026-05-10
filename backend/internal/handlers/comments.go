package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentsHandler struct {
	pool *pgxpool.Pool
}

func NewCommentsHandler(pool *pgxpool.Pool) *CommentsHandler {
	return &CommentsHandler{pool: pool}
}

func (h *CommentsHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
