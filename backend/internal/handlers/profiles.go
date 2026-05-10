package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfilesHandler struct {
	pool *pgxpool.Pool
}

func NewProfilesHandler(pool *pgxpool.Pool) *ProfilesHandler {
	return &ProfilesHandler{pool: pool}
}

func (h *ProfilesHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
