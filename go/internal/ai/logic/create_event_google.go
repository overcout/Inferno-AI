package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/overcout/Inferno-AI/internal/ai/engine"
	"github.com/overcout/Inferno-AI/internal/logger"
)

// CreateEventGoogleCommand holds data for action_create_event_google
type CreateEventGoogleCommand struct {
	Action          string `json:"action"`
	Title           string `json:"title"`
	StartTime       string `json:"start_time"`
	DurationMinutes int    `json:"duration_minutes"`
}

func (c *CreateEventGoogleCommand) RenderMessage() string {
	return fmt.Sprintf(
		"📅 Create event:\n• Title: %s\n• Start: %s\n• Duration: %d minutes",
		c.Title, c.StartTime, c.DurationMinutes,
	)
}

// GenerateCreateEventGoogle parses AI response into a CreateEventGoogleCommand
func GenerateCreateEventGoogle(ai engine.AIEngine, userPrompt string) (*CreateEventGoogleCommand, error) {
	logger.Info.Println("Generating action_create_event_google payload")

	today := time.Now().Format("2006-01-02")
	systemPrompt := "" +
		"You are a command compiler. Today is " + today + ". " +
		"Convert every user instruction into strict JSON format like:\n" +
		"{\n" +
		"  \"action\": \"action_create_event_google\",\n" +
		"  \"title\": \"Meeting with Ivan\",\n" +
		"  \"start_time\": \"2025-03-22T15:00:00\",\n" +
		"  \"duration_minutes\": 30\n" +
		"}\n" +
		"Respond with JSON only. NEVER include explanations, comments, Markdown blocks, or any additional text. Always return raw JSON."

	prompt := systemPrompt + "\n\nUser: " + userPrompt

	raw, err := ai.ProcessPrompt(prompt)
	if err != nil {
		return nil, err
	}

	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.Trim(cleaned, "`")

	re := regexp.MustCompile("(?s)^```[a-zA-Z]*\\n(.*)```$")
	if matches := re.FindStringSubmatch(cleaned); len(matches) == 2 {
		cleaned = matches[1]
	}

	var cmd CreateEventGoogleCommand
	if err := json.Unmarshal([]byte(cleaned), &cmd); err != nil {
		logger.Error.Println("Invalid JSON returned by model:", cleaned)
		return nil, errors.New("model returned invalid JSON: " + cleaned)
	}

	logger.Info.Println("Command successfully parsed:", cmd)
	return &cmd, nil
}
