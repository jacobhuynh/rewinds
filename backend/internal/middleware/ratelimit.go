package middleware

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
}

type visitor struct {
	count   int
	resetAt time.Time
}

var limiter = &rateLimiter{visitors: make(map[string]*visitor)}

// RateLimit allows up to limit requests per window per IP.
func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			limiter.mu.Lock()
			v, ok := limiter.visitors[ip]
			if !ok || time.Now().After(v.resetAt) {
				limiter.visitors[ip] = &visitor{count: 1, resetAt: time.Now().Add(window)}
				limiter.mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}
			v.count++
			if v.count > limit {
				limiter.mu.Unlock()
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			limiter.mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}
