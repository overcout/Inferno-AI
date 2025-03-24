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
	Title           string `json:"title"`
	StartTime       string `json:"start_time"`
	DurationMinutes int    `json:"duration_minutes"`
}

// ListEventsGoogleCommand holds the parsed parameters for querying events
type ListEventsGoogleCommand struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (c *CreateEventGoogleCommand) RenderMessage() string {
	return fmt.Sprintf(
		"📅 Create event:\n• Title: %s\n• Start: %s\n• Duration: %d minutes",
		c.Title, c.StartTime, c.DurationMinutes,
	)
}

func (c *ListEventsGoogleCommand) RenderMessage() string {
	return fmt.Sprintf(
		"📆 Show events from %s to %s",
		c.From, c.To,
	)
}

func isValidDate(value string) bool {
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}

func isValidDateTime(value string) bool {
	_, err := time.Parse("2006-01-02T15:04:05", value)
	return err == nil
}

// GenerateCreateEventGoogle parses AI response into a CreateEventGoogleCommand
func GenerateCreateEventGoogle(ai engine.AIEngine, userPrompt string) (*CreateEventGoogleCommand, error) {
	logger.Info.Println("Generating action_create_event_google payload")

	today := time.Now().Format("2006-01-02")
	systemPrompt := "" +
		"You are a command compiler. Today is " + today + ". " +
		"You MUST return valid JSON matching the example exactly. " +
		"Example:\n" +
		"{\n" +
		"  \"title\": \"Meeting with Ivan\",\n" +
		"  \"start_time\": \"2025-03-22T15:00:00\",\n" +
		"  \"duration_minutes\": 30\n" +
		"}\n" +
		"Respond with JSON only. NEVER include explanations, comments, Markdown blocks, or any additional text."

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

	if !isValidDateTime(cmd.StartTime) {
		return nil, errors.New("invalid or malformed command: " + cleaned)
	}

	logger.Info.Println("Command successfully parsed:", cmd)
	return &cmd, nil
}

// GenerateListEventsGoogle parses the user's natural language into a structured command
func GenerateListEventsGoogle(ai engine.AIEngine, userPrompt string) (*ListEventsGoogleCommand, error) {
	today := time.Now().Format("2006-01-02")
	logger.Info.Println("Generating action_list_events_google payload")

	systemPrompt := "" +
		"You are a command compiler. Today is " + today + ". " +
		"You MUST return valid JSON matching the example exactly. " +
		"Example:\n" +
		"{\n" +
		"  \"from\": \"2025-03-22\",\n" +
		"  \"to\": \"2025-03-24\"\n" +
		"}\n" +
		"Respond with JSON only. NEVER include explanations, comments or any other formatting."

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

	var cmd ListEventsGoogleCommand
	if err := json.Unmarshal([]byte(cleaned), &cmd); err != nil {
		logger.Error.Println("Invalid JSON returned by model:", cleaned)
		return nil, errors.New("model returned invalid JSON: " + cleaned)
	}

	if !isValidDate(cmd.From) || !isValidDate(cmd.To) {
		return nil, errors.New("invalid or malformed command: " + cleaned)
	}

	logger.Info.Println("Parsed command:", cmd)
	return &cmd, nil
}
