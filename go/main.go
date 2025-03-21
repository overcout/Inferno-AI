package main

import (
	"fmt"
	"github.com/overcout/Inferno-AI/internal/ai"
	"github.com/overcout/Inferno-AI/internal/ai/providers"
	"github.com/overcout/Inferno-AI/internal/logger"
)

func main() {
	// Initialize logger
	logger.InitConsole()
	logger.Info.Println("Bot started")

	// Create AI engine (Ollama implementation)
	ollama := providers.NewOllamaEngine("http://localhost:11434", "phi4")
	controller := ai.NewAIController(ollama)

	// Example prompt
	prompt := "Make a meeting with Ivan at 15:30 tomorrow"

	cmd, err := controller.ProcessRequest(prompt)
	if err != nil {
		logger.Error.Println("Error processing request:", err)
		return
	}

	if cmd == nil {
		logger.Warning.Println("No command returned. Possibly unrecognized action.")
		return
	}

	fmt.Println("\nParsed command:")
	fmt.Printf("Action: %s\nTitle: %s\nStart: %s\nDuration: %d minutes\n",
		cmd.Action, cmd.Title, cmd.StartTime, cmd.DurationMinutes)
}
