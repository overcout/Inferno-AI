package ai

import (
	"fmt"
	"time"

	"github.com/overcout/Inferno-AI/internal/ai/actions"
	"github.com/overcout/Inferno-AI/internal/ai/engine"
	"github.com/overcout/Inferno-AI/internal/ai/logic"
	"github.com/overcout/Inferno-AI/internal/store"
)

type AIController struct {
	Engine        engine.AIEngine
	Store         *store.Store
	CurrentUserID int64
}

func NewAIController(engine engine.AIEngine, store *store.Store, userID int64) *AIController {
	return &AIController{
		Engine:        engine,
		Store:         store,
		CurrentUserID: userID,
	}
}

func (c *AIController) ProcessRequest(prompt string) (logic.Renderable, error) {
	action, err := actions.DetectAction(c.Engine, prompt)
	if err != nil {
		return nil, err
	}

	user, err := c.Store.GetOrCreateUser(c.CurrentUserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	switch action {

	case "action_create_event_google":
		cmd, err := logic.GenerateCreateEventGoogle(c.Engine, prompt)
		if err != nil {
			return nil, err
		}
		err = logic.CreateEventGoogle(user, cmd)
		if err != nil {
			return nil, fmt.Errorf("event creation failed: %w", err)
		}
		return cmd, nil

	case "action_list_events_google":
		cmd, err := logic.GenerateListEventsGoogle(c.Engine, prompt)
		if err != nil {
			return nil, err
		}

		from, err := time.Parse("2006-01-02", cmd.From)
		if err != nil {
			return nil, fmt.Errorf("invalid 'from' format: %w", err)
		}
		to, err := time.Parse("2006-01-02", cmd.To)
		if err != nil {
			return nil, fmt.Errorf("invalid 'to' format: %w", err)
		}

		events, err := logic.ListEventsGoogle(user, from, to)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch events: %w", err)
		}

		return logic.RenderText(logic.RenderEvents(events)), nil

	case "action_undefined":
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}
