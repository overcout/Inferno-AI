package ai

import (
	"fmt"

	"github.com/overcout/Inferno-AI/internal/ai/actions"
	"github.com/overcout/Inferno-AI/internal/ai/engine"
	"github.com/overcout/Inferno-AI/internal/ai/logic"
	"github.com/overcout/Inferno-AI/internal/store"
)

// AIController orchestrates AI requests and action execution
type AIController struct {
	Engine        engine.AIEngine
	Store         *store.Store
	CurrentUserID int64 // Telegram user ID
}

// NewAIController creates a new controller with injected dependencies
func NewAIController(engine engine.AIEngine, store *store.Store, userID int64) *AIController {
	return &AIController{
		Engine:        engine,
		Store:         store,
		CurrentUserID: userID,
	}
}

// ProcessRequest determines the action, builds command, executes logic
func (c *AIController) ProcessRequest(prompt string) (logic.Renderable, error) {
	action, err := actions.DetectAction(c.Engine, prompt)
	if err != nil {
		return nil, err
	}

	switch action {
	case "action_create_event_google":
		cmd, err := logic.GenerateCreateEventGoogle(c.Engine, prompt)
		if err != nil {
			return nil, err
		}

		user, err := c.Store.GetOrCreateUser(c.CurrentUserID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}

		err = logic.CreateEventGoogle(user, cmd)
		if err != nil {
			return nil, fmt.Errorf("event creation failed: %w", err)
		}

		return cmd, nil

	case "action_undefined":
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}
