package main

import (
	"github.com/overcout/Inferno-AI/internal/ai"
	"github.com/overcout/Inferno-AI/internal/ai/providers"
	"github.com/overcout/Inferno-AI/internal/config"
	"github.com/overcout/Inferno-AI/internal/logger"
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

	var controller *ai.AIController

	switch cfg.EngineType {
	case "ollama":
		engine := providers.NewOllamaEngine(cfg.OllamaURL, cfg.OllamaModel)
		controller = ai.NewAIController(engine)
	default:
		logger.Error.Println("Unknown engine type:", cfg.EngineType)
		return
	}

	// 🔌 Start Telegram bot if token is present
	if cfg.TelegramToken != "" {
		telegram.StartTelegramBot(cfg.TelegramToken, controller)
	}
}
