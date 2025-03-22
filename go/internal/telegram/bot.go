package telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/overcout/Inferno-AI/internal/ai"
	"github.com/overcout/Inferno-AI/internal/config"
	"github.com/overcout/Inferno-AI/internal/store"
	"github.com/overcout/Inferno-AI/internal/tools"
)

// StartTelegramBot launches the bot and listens for updates
func StartTelegramBot(token string, controller *ai.AIController, db *store.Store, cfg *config.Config) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to start Telegram bot: %v", err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				bot.Send(tgbotapi.NewMessage(chatID, "👋 Welcome! Send me a request in natural language, and I'll help you!"))

			case "help":
				bot.Send(tgbotapi.NewMessage(chatID, "💡 Example: 'Create a meeting with Anna tomorrow at 3 PM for 1 hour'"))

			case "auth":
				token := tools.GenerateToken(32)
				_, err := db.CreateAuthLink(token, int64(userID), 10*60)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "❌ Failed to create auth link: "+err.Error()))
					continue
				}
				link := fmt.Sprintf("%s/oauth?token=%s", cfg.OAuthPublicURL, token)
				bot.Send(tgbotapi.NewMessage(chatID, "🔐 Authorize your Google account:\n"+link))

			default:
				bot.Send(tgbotapi.NewMessage(chatID, "Unknown command. Try /help"))
			}
			continue
		}

		prompt := strings.TrimSpace(update.Message.Text)
		if prompt == "" {
			bot.Send(tgbotapi.NewMessage(chatID, "Please send text message only."))
			continue
		}

		userController := ai.NewAIController(controller.Engine, controller.Store, int64(userID))

		result, err := userController.ProcessRequest(prompt)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Error: "+err.Error()))
			continue
		}
		if result == nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Sorry, I didn't understand your request."))
			continue
		}

		bot.Send(tgbotapi.NewMessage(chatID, result.RenderMessage()))
	}
}
