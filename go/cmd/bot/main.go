package main

import (
	"github.com/overcout/Inferno-AI/internal/ai"
	"github.com/overcout/Inferno-AI/internal/ai/providers"
	"github.com/overcout/Inferno-AI/internal/config"
	"github.com/overcout/Inferno-AI/internal/logger"
	"github.com/overcout/Inferno-AI/internal/store"
	"github.com/overcout/Inferno-AI/internal/telegram"
)

func main() {
	cfg := config.LoadConfig()

	switch cfg.LogMode {
	case "console":
		logger.InitConsole()
	default:
		logger.Init()
	}
	logger.Info.Println("Bot started")

	// Init store (MariaDB via GORM)
	if cfg.DBDSN == "" {
		logger.Error.Println("DB_DSN not set")
		return
	}

	storeInstance, err := store.NewStore(cfg.DBDSN)
	if err != nil {
		logger.Error.Println("Failed to connect to DB:", err)
		return
	}

	if err := storeInstance.InitSchema(); err != nil {
		logger.Error.Println("Failed to initialize schema:", err)
		return
	}

	// Init AI controller
	var controller *ai.AIController
	switch cfg.EngineType {
	case "ollama":
		engine := providers.NewOllamaEngine(cfg.OllamaURL, cfg.OllamaModel)
		controller = ai.NewAIController(engine, storeInstance, 0)
	default:
		logger.Error.Println("Unknown engine type:", cfg.EngineType)
		return
	}

	// 🤖 Start Telegram bot
	if cfg.TelegramToken != "" {
		telegram.StartTelegramBot(cfg.TelegramToken, controller, storeInstance, cfg)
	}
}
