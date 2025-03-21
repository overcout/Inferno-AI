package ai

import (
	"errors"

	"github.com/overcout/Inferno-AI/internal/ai/actions"
	"github.com/overcout/Inferno-AI/internal/ai/engine"
	"github.com/overcout/Inferno-AI/internal/ai/logic"
	"github.com/overcout/Inferno-AI/internal/logger"
)

// AIController coordinates between action detection and specific logic
type AIController struct {
	Engine engine.AIEngine
}

func NewAIController(e engine.AIEngine) *AIController {
	return &AIController{Engine: e}
}

// ProcessRequest handles a prompt and returns a Renderable command
func (c *AIController) ProcessRequest(prompt string) (logic.Renderable, error) {
	action, err := actions.DetectAction(c.Engine, prompt)
	if err != nil {
		return nil, err
	}

	switch action {
	case "action_create_event_google":
		return logic.GenerateCreateEvent(c.Engine, prompt)
	case "action_send_email_google":
		return logic.GenerateSendEmail(c.Engine, prompt)
	case "action_undefined":
		logger.Warning.Println("AI could not identify a valid action")
		return nil, nil
	default:
		logger.Error.Println("Unknown or unsupported action:", action)
		return nil, errors.New("unhandled action: " + action)
	}
}
