package actions

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/overcout/Inferno-AI/internal/ai/engine"
	"github.com/overcout/Inferno-AI/internal/logger"
)

// ActionDetectionResult represents the result of intent classification.
type ActionDetectionResult struct {
	Action string `json:"action"`
}

// List of supported actions with descriptions.
var SupportedActions = []struct {
	Key         string
	Description string
}{
	{"action_create_event", "Create a calendar event in Google Calendar"},
	{"action_send_email_gmail", "Send an email using Gmail API"},
	{"action_list_events", "List today's calendar events"},
	{"action_undefined", "Use this when user intent doesn't match any known action"},
}

// DetectAction uses the AI engine to classify the user's intent and return an action.
func DetectAction(engine engine.AIEngine, userPrompt string) (string, error) {
	logger.Info.Println("Detecting action from prompt")

	var b strings.Builder
	b.WriteString("You are an intent classifier. Your task is to choose one action key from the list below that best describes the user's request.\n")
	b.WriteString("Respond with JSON: {\"action\": \"...\"} and nothing else.\n\n")
	b.WriteString("Available actions:\n")
	for _, a := range SupportedActions {
		b.WriteString("- " + a.Key + ": " + a.Description + "\n")
	}
	b.WriteString("\nUser request: " + userPrompt)

	raw, err := engine.ProcessPrompt(b.String())
	if err != nil {
		return "", err
	}

	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.Trim(cleaned, "`")

	re := regexp.MustCompile("(?s)^```[a-zA-Z]*\\n(.*)```$")
	if matches := re.FindStringSubmatch(cleaned); len(matches) == 2 {
		cleaned = matches[1]
	}

	var result ActionDetectionResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		logger.Error.Println("Failed to parse action JSON:", cleaned)
		return "", errors.New("invalid JSON for action: " + cleaned)
	}

	logger.Info.Println("Detected action:", result.Action)
	return result.Action, nil
}
