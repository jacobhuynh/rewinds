package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ComparisonsHandler struct {
	pool  *pgxpool.Pool
	redis *redis.Client
}

func NewComparisonsHandler(pool *pgxpool.Pool, redis *redis.Client) *ComparisonsHandler {
	return &ComparisonsHandler{pool: pool, redis: redis}
}

func (h *ComparisonsHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
