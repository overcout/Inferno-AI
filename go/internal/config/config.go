package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	EngineType    string
	OllamaURL     string
	OllamaModel   string
	LogMode       string
	TelegramToken string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{}

	flag.StringVar(&cfg.EngineType, "engine", getEnv("AI_ENGINE", "ollama"), "AI engine to use")
	flag.StringVar(&cfg.OllamaURL, "ollama-url", getEnv("OLLAMA_URL", "http://localhost:11434"), "Ollama URL")
	flag.StringVar(&cfg.OllamaModel, "ollama-model", getEnv("OLLAMA_MODEL", "phi4"), "Ollama model name")
	flag.StringVar(&cfg.LogMode, "log-mode", getEnv("LOG_MODE", "console"), "Logging mode")
	flag.StringVar(&cfg.TelegramToken, "telegram-token", getEnv("TELEGRAM_TOKEN", ""), "Telegram bot token")

	flag.Parse()
	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
