package actions

import (
	"strings"

	"github.com/overcout/Inferno-AI/internal/ai/engine"
	"github.com/overcout/Inferno-AI/internal/logger"
)

type ActionDetectionResult struct {
	Action string `json:"action"`
}

var SupportedActions = []struct {
	Key         string
	Description string
}{
	{"action_create_event_google", "Create a calendar event in Google Calendar. Use this action when the user requests to schedule a new event with specific details like time, date, and event description."},
	{"action_list_events_google", "List events from Google Calendar. This action is triggered when the user wants to see a list of upcoming events in their calendar. This action will not handle questions like 'What is the closest event?'."},
	{"action_undefined", "Use this when user intent doesn't match any known action. This is used when the user's request doesn't match any available actions. It ensures that unrecognized or invalid inputs don't trigger a system error."},
}

func DetectAction(engine engine.AIEngine, userPrompt string) (string, error) {
	logger.Info.Println("Detecting action from prompt")

	var b strings.Builder
	b.WriteString("You are an intent classifier. Your task is to choose one action key from the list below that best describes the user's request.\n")
	b.WriteString("Respond with only the action key (e.g., 'action_create_event_google') and nothing else.\n\n")
	b.WriteString("Available actions:\n")
	for _, a := range SupportedActions {
		b.WriteString("- " + a.Key + ": " + a.Description + "\n")
	}
	b.WriteString("\nUser request: " + userPrompt)

	logger.Info.Println("Prepared prompt for AI engine:", b.String())

	raw, err := engine.ProcessPrompt(b.String())
	if err != nil {
		logger.Error.Println("Error during AI processing:", err)
		return "action_undefined", err
	}

	logger.Info.Println("AI response:", raw)

	cleaned := strings.TrimSpace(raw)

	validAction := false
	for _, a := range SupportedActions {
		if a.Key == cleaned {
			validAction = true
			break
		}
	}

	if !validAction {
		logger.Error.Println("Invalid action detected. Returning 'action_undefined'.")
		return "action_undefined", nil
	}

	logger.Info.Println("Detected valid action:", cleaned)
	return cleaned, nil
}
