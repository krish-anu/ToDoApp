package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var ErrMissingMongoURI = errors.New("MONGODB_URI not set")

type Config struct {
	MongoURI                   string
	DatabaseName               string
	CollectionName             string
	ReminderCollectionName     string
	WorkflowRunsCollectionName string
	Port                       string
	AllowOrigins               string
	OpenAIAPIKey               string
	OpenAIModel                string
	OpenAIBaseURL              string
	WorkflowDefaultUserID      string
	WorkflowDefaultTimezone    string
}

func Load() (Config, error) {
	// Load local environment values when available (safe in production if files are absent).
	_ = godotenv.Load(".env", "backend/.env")

	cfg := Config{
		MongoURI:                   os.Getenv("MONGODB_URI"),
		DatabaseName:               getEnv("MONGODB_DB", "golang_db"),
		CollectionName:             getEnv("MONGODB_COLLECTION", "todos"),
		ReminderCollectionName:     getEnv("MONGODB_REMINDER_COLLECTION", "reminders"),
		WorkflowRunsCollectionName: getEnv("MONGODB_WORKFLOW_RUNS_COLLECTION", "workflow_runs"),
		Port:                       getEnv("PORT", "5000"),
		AllowOrigins:               getEnv("ALLOW_ORIGINS", "http://localhost:5173,http://localhost:5174"),
		OpenAIAPIKey:               os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:                getEnv("OPENAI_MODEL", "gpt-4.1-mini"),
		OpenAIBaseURL:              getEnv("OPENAI_BASE_URL", "https://api.openai.com"),
		WorkflowDefaultUserID:      getEnv("WORKFLOW_DEFAULT_USER_ID", "local-user"),
		WorkflowDefaultTimezone:    getEnv("WORKFLOW_DEFAULT_TIMEZONE", "UTC"),
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
