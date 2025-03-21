package telegram

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/overcout/Inferno-AI/internal/ai"
)

// StartTelegramBot launches the bot and listens for updates
func StartTelegramBot(token string, controller *ai.AIController) {
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

		// Handle /start and /help commands
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				bot.Send(tgbotapi.NewMessage(chatID, "👋 Welcome! Send me a request in natural language, and I'll help you!"))
			case "help":
				bot.Send(tgbotapi.NewMessage(chatID, "💡 Example: 'Create a meeting with Anna tomorrow at 3 PM for 1 hour'"))
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

		result, err := controller.ProcessRequest(prompt)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Error: "+err.Error()))
			continue
		}
		if result == nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Sorry, I didn't understand your request."))
			continue
		}

		// ✅ Use Renderable interface
		message := result.RenderMessage()
		bot.Send(tgbotapi.NewMessage(chatID, message))
	}
}
