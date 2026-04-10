package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var ErrMissingMongoURI = errors.New("MONGODB_URI not set")

type Config struct {
	MongoURI       string
	DatabaseName   string
	CollectionName string
	Port           string
	AllowOrigins   string
}

func Load() (Config, error) {
	// Load local environment values when available (safe in production if file is absent).
	_ = godotenv.Load(".env")

	cfg := Config{
		MongoURI:       os.Getenv("MONGODB_URI"),
		DatabaseName:   getEnv("MONGODB_DB", "golang_db"),
		CollectionName: getEnv("MONGODB_COLLECTION", "todos"),
		Port:           getEnv("PORT", "5000"),
		AllowOrigins:   getEnv("ALLOW_ORIGINS", "http://localhost:5173,http://localhost:5174"),
	}

	if cfg.MongoURI == "" {
		return Config{}, ErrMissingMongoURI
	}

	return cfg, nil
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
