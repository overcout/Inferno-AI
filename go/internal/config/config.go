package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configurable options for the app
type Config struct {
	EngineType          string
	OllamaURL           string
	OllamaModel         string
	LogMode             string
	TelegramToken       string
	GoogleClientID      string
	GoogleClientSecret  string
	OAuthRedirectURL    string
	OAuthPublicURL      string
	DBDSN               string
}

// LoadConfig parses flags, env vars, and loads .env file if present
func LoadConfig() *Config {
	_ = godotenv.Load() // load from .env if exists (safe to ignore if missing)

	cfg := &Config{}

	flag.StringVar(&cfg.EngineType, "engine", getEnv("AI_ENGINE", "ollama"), "AI engine to use (e.g., ollama)")
	flag.StringVar(&cfg.OllamaURL, "ollama-url", getEnv("OLLAMA_URL", "http://localhost:11434"), "URL of the Ollama service")
	flag.StringVar(&cfg.OllamaModel, "ollama-model", getEnv("OLLAMA_MODEL", "phi4"), "Model name for Ollama")
	flag.StringVar(&cfg.LogMode, "log-mode", getEnv("LOG_MODE", "console"), "Logging mode: console or file")
	flag.StringVar(&cfg.TelegramToken, "telegram-token", getEnv("TELEGRAM_TOKEN", ""), "Telegram bot token")
	flag.StringVar(&cfg.GoogleClientID, "google-client-id", getEnv("GOOGLE_CLIENT_ID", ""), "Google OAuth client ID")
	flag.StringVar(&cfg.GoogleClientSecret, "google-client-secret", getEnv("GOOGLE_CLIENT_SECRET", ""), "Google OAuth client secret")
	flag.StringVar(&cfg.OAuthRedirectURL, "oauth-redirect-url", getEnv("OAUTH_REDIRECT_URL", "http://localhost:8080/oauth/callback"), "OAuth redirect URL")
	flag.StringVar(&cfg.OAuthPublicURL, "oauth-public-url", getEnv("OAUTH_PUBLIC_URL", "https://your-auth-url.com"), "OAuth public base URL used for Telegram links")
	flag.StringVar(&cfg.DBDSN, "db-dsn", getEnv("DB_DSN", ""), "Database DSN for MariaDB")

	flag.Parse()

	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}