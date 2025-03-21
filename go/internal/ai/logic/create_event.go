package logic

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/overcout/Inferno-AI/internal/ai/engine"
	"github.com/overcout/Inferno-AI/internal/logger"
)

// Command represents the parsed command structure.
type Command struct {
	Action          string `json:"action"`
	Title           string `json:"title"`
	StartTime       string `json:"start_time"`
	DurationMinutes int    `json:"duration_minutes"`
}

// GenerateCreateEvent parses prompt into structured command for calendar creation.
func GenerateCreateEvent(ai engine.AIEngine, userPrompt string) (*Command, error) {
	logger.Info.Println("Generating create_event payload")

	today := time.Now().Format("2006-01-02")
	systemPrompt := "" +
		"You are a command compiler. Today is " + today + ". " +
		"Convert every user instruction into strict JSON format like:\n" +
		"{\n" +
		"  \"action\": \"create_event\",\n" +
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

	var cmd Command
	if err := json.Unmarshal([]byte(cleaned), &cmd); err != nil {
		logger.Error.Println("Invalid JSON returned by model:", cleaned)
		return nil, errors.New("model returned invalid JSON: " + cleaned)
	}

	logger.Info.Println("Command successfully parsed:", cmd)
	return &cmd, nil
}
