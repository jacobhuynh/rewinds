package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jacobhuynh/rewinds/backend/config"
	"github.com/jacobhuynh/rewinds/backend/internal/db"
	"github.com/jacobhuynh/rewinds/backend/internal/handlers"
	custommiddleware "github.com/jacobhuynh/rewinds/backend/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env from repo root (one level up from /backend).
	// Missing file is fine (Railway injects env vars directly);
	// any other error (permissions, malformed) is a real warning.
	if err := godotenv.Load("../.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("warning: could not load .env: %v", err)
	}

	cfg := config.Load()
	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()
	log.Println("database connected")

	redisClient, err := db.NewRedis(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("redis connected")

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(custommiddleware.RateLimit(100, time.Minute))

	// Public routes
	r.Get("/health", healthHandler)

	// Auth
	authHandler := handlers.NewAuthHandler(pool, cfg)
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	spotifyAuth := handlers.NewSpotifyHandler(pool, cfg, redisClient)
	r.Post("/auth/spotify", spotifyAuth.Exchange)
	r.Get("/auth/spotify/callback", spotifyAuth.Callback)
	r.Post("/auth/spotify/refresh", spotifyAuth.Refresh)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(custommiddleware.Auth(cfg.JWTSecret))

		artistsH := handlers.NewArtistsHandler(pool)
		r.Get("/artists/{id}", artistsH.NotImplemented)
		r.Get("/artists/matchup", artistsH.NotImplemented)

		albumsH := handlers.NewAlbumsHandler(pool)
		r.Get("/albums/{id}", albumsH.NotImplemented)

		tracksH := handlers.NewTracksHandler(pool)
		r.Get("/tracks/{id}", tracksH.NotImplemented)

		votesH := handlers.NewVotesHandler(pool)
		r.Post("/votes", votesH.NotImplemented)

		ratingsH := handlers.NewRatingsHandler(pool, redisClient)
		r.Post("/ratings", ratingsH.NotImplemented)
		r.Get("/ratings", ratingsH.NotImplemented)
		r.Put("/ratings/{id}", ratingsH.NotImplemented)
		r.Delete("/ratings/{id}", ratingsH.NotImplemented)

		comparisonsH := handlers.NewComparisonsHandler(pool, redisClient)
		r.Post("/comparisons", comparisonsH.NotImplemented)
		r.Get("/comparisons/next", comparisonsH.NotImplemented)

		leaderboardH := handlers.NewLeaderboardHandler(pool)
		r.Get("/leaderboard", leaderboardH.NotImplemented)
		r.Get("/leaderboard/rising", leaderboardH.NotImplemented)
		r.Get("/leaderboard/users", leaderboardH.NotImplemented)

		profilesH := handlers.NewProfilesHandler(pool)
		r.Get("/profiles/{username}", profilesH.NotImplemented)
		r.Put("/profiles/{username}", profilesH.NotImplemented)
		r.Post("/profiles/avatar", profilesH.NotImplemented)

		predictionsH := handlers.NewPredictionsHandler(pool)
		r.Post("/predictions", predictionsH.NotImplemented)
		r.Get("/predictions", predictionsH.NotImplemented)
		r.Get("/predictions/available", predictionsH.NotImplemented)

		pointsH := handlers.NewPointsHandler(pool)
		r.Get("/users/{id}/rank", pointsH.NotImplemented)

		commentsH := handlers.NewCommentsHandler(pool)
		r.Post("/comments", commentsH.NotImplemented)
		r.Get("/comments", commentsH.NotImplemented)
		r.Delete("/comments/{id}", commentsH.NotImplemented)

		foldersH := handlers.NewFoldersHandler(pool)
		r.Post("/folders", foldersH.NotImplemented)
		r.Get("/folders", foldersH.NotImplemented)
		r.Get("/folders/{id}", foldersH.NotImplemented)
		r.Put("/folders/{id}", foldersH.NotImplemented)
		r.Delete("/folders/{id}", foldersH.NotImplemented)
		r.Post("/folders/{id}/ratings", foldersH.NotImplemented)
		r.Delete("/folders/{id}/ratings/{rating_id}", foldersH.NotImplemented)

		questsH := handlers.NewQuestsHandler(pool)
		r.Get("/quests", questsH.NotImplemented)
		r.Post("/quests/{id}/claim", questsH.NotImplemented)

		promoH := handlers.NewPromoHandler(pool)
		r.Get("/promo/available", promoH.NotImplemented)
		r.Post("/promo/{campaign_id}/rate", promoH.NotImplemented)

		toursH := handlers.NewToursHandler(pool)
		r.Post("/tour-alerts", toursH.NotImplemented)
		r.Get("/tour-alerts", toursH.NotImplemented)

		spotifyH := handlers.NewSpotifyHandler(pool, cfg, redisClient)
		r.Get("/spotify/onboarding", spotifyH.Onboarding)
		r.Get("/spotify/playlists", spotifyH.Playlists)
		r.Get("/spotify/playlists/{playlist_id}/tracks", spotifyH.PlaylistTracks)

		discoverH := handlers.NewArtistsHandler(pool)
		r.Get("/discover", discoverH.NotImplemented)
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second, // matches chi Timeout middleware
		IdleTimeout:  120 * time.Second,
	}

	// Start server in background so main goroutine can listen for shutdown signals.
	go func() {
		log.Printf("server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	// Block until SIGTERM (Railway deploy) or SIGINT (Ctrl+C) received.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("shutdown signal received, draining connections...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
	log.Println("server stopped")
	// Deferred pool.Close() and redisClient.Close() now run cleanly.
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
