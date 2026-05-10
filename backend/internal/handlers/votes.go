package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VotesHandler struct {
	pool *pgxpool.Pool
}

func NewVotesHandler(pool *pgxpool.Pool) *VotesHandler {
	return &VotesHandler{pool: pool}
}

func (h *VotesHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
