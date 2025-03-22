# Inferno-AI

Inferno-AI is a Telegram bot that integrates AI and Google services. It allows you to create events in Google Calendar using natural language through Telegram.

## Architecture
Modules:
- **Telegram Bot** - accepts commands and interacts with the user.
- **AI Controller** - processes messages and determines actions using AI.
- **OAuth Service** - a separate HTTP service for authorization through Google.
- **Store** - MariaDB database (with tokens and users).

The scheme is in the attached file `architecture.png`.

## Launch
1. Install Go 1.19+, MariaDB, Telegram Bot and Google OAuth2 client.
2. Configure `.env` in `go/`.
3. Start OAuth: `go run cmd/oauth/main.go`.
4. Start the bot: `go run cmd/bot/main.go`.

Example request: _"Create a meeting with Anna tomorrow at 3:00 PM for 1 hour"._

## .env example
```
TELEGRAM_TOKEN=<token>
DB_DSN=<dsn>
GOOGLE_CLIENT_ID=<id>
GOOGLE_CLIENT_SECRET=<secret>
OAUTH_PUBLIC_URL=https://<url>
OAUTH_REDIRECT_URL=https://<url>/oauth/callback
AI_ENGINE=ollama
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=phi4
LOG_MODE=console
```