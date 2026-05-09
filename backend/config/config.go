package config

import (
	"log"
	"os"
)

type Config struct {
	Port                string
	DatabaseURL         string
	RedisURL            string
	JWTSecret           string
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURI  string
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8080"),
		DatabaseURL:         mustEnv("SUPABASE_DB_URL"),
		RedisURL:            getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:           mustEnv("JWT_SECRET"),
		SpotifyClientID:     mustEnv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: mustEnv("SPOTIFY_CLIENT_SECRET"),
		SpotifyRedirectURI:  mustEnv("SPOTIFY_REDIRECT_URI"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}
