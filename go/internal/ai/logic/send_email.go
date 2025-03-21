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

// SendEmailCommand represents action_action_send_email_google
type SendEmailCommand struct {
	Action  string `json:"action"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (c *SendEmailCommand) RenderMessage() string {
	return fmt.Sprintf(
		"📧 Send Email:\n• To: %s\n• Subject: %s\n• Body:\n%s",
		c.To, c.Subject, c.Body,
	)
}

// GenerateSendEmail parses AI response into SendEmailCommand
func GenerateSendEmail(ai engine.AIEngine, userPrompt string) (*SendEmailCommand, error) {
	logger.Info.Println("Generating action_send_email_google payload")

	today := time.Now().Format("2006-01-02")
	systemPrompt := "" +
		"You are a command compiler. Today is " + today + ". " +
		"Convert the user's request into a JSON to send email:\n" +
		"{\n" +
		"  \"action\": \"action_send_email_google\",\n" +
		"  \"to\": \"john@example.com\",\n" +
		"  \"subject\": \"Job Offer\",\n" +
		"  \"body\": \"Hey John, I'm hiring you!\"\n" +
		"}\n" +
		"Respond ONLY with raw JSON."

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

	var cmd SendEmailCommand
	if err := json.Unmarshal([]byte(cleaned), &cmd); err != nil {
		logger.Error.Println("Invalid JSON returned by model:", cleaned)
		return nil, errors.New("model returned invalid JSON: " + cleaned)
	}

	logger.Info.Println("Email command successfully parsed:", cmd)
	return &cmd, nil
}
