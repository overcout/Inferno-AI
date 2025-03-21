package ai

import (
	"errors"
	"github.com/overcout/Inferno-AI/internal/ai/actions"
	"github.com/overcout/Inferno-AI/internal/ai/engine"
	"github.com/overcout/Inferno-AI/internal/ai/logic"
	"github.com/overcout/Inferno-AI/internal/logger"
)

// AIController provides a unified interface to detect action and generate payload.
type AIController struct {
	Engine engine.AIEngine
}

// NewAIController creates a new controller.
func NewAIController(e engine.AIEngine) *AIController {
	return &AIController{Engine: e}
}

// ProcessRequest splits request into 2 steps: action detection, then payload generation.
func (c *AIController) ProcessRequest(prompt string) (*logic.Command, error) {
	// Step 1: Detect action
	action, err := actions.DetectAction(c.Engine, prompt)
	if err != nil {
		return nil, err
	}

	// Step 2: Generate payload based on action
	switch action {
	case "action_create_event":
		return logic.GenerateCreateEvent(c.Engine, prompt)
	case "action_undefined":
		logger.Warning.Println("AI could not identify a valid action")
		return nil, nil
	default:
		logger.Error.Println("Unknown or unsupported action:", action)
		return nil, errors.New("unhandled action: " + action)
	}
}