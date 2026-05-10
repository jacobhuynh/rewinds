package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FoldersHandler struct {
	pool *pgxpool.Pool
}

func NewFoldersHandler(pool *pgxpool.Pool) *FoldersHandler {
	return &FoldersHandler{pool: pool}
}

func (h *FoldersHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
