package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type RatingsHandler struct {
	pool  *pgxpool.Pool
	redis *redis.Client
}

func NewRatingsHandler(pool *pgxpool.Pool, redis *redis.Client) *RatingsHandler {
	return &RatingsHandler{pool: pool, redis: redis}
}

func (h *RatingsHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
