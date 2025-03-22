package actions_test

import (
	"testing"

	"github.com/overcout/Inferno-AI/internal/ai/actions"
	"github.com/overcout/Inferno-AI/internal/logger"
)

type mockAI struct {
	Response string
	Err      error
}

func (m *mockAI) ProcessPrompt(prompt string) (string, error) {
	return m.Response, m.Err
}

func TestDetectAction_Undefined(t *testing.T) {
	logger.InitConsole()
	mock := &mockAI{
		Response: `{"action": "action_undefined"}`,
	}
	action, err := actions.DetectAction(mock, "Tell me a joke")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action != "action_undefined" {
		t.Errorf("expected action_undefined, got %s", action)
	}
}

func TestDetectAction_CreateEventGoogle(t *testing.T) {
	logger.InitConsole()
	mock := &mockAI{
		Response: `{"action": "action_create_event_google"}`,
	}
	action, err := actions.DetectAction(mock, "Schedule a meeting")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action != "action_create_event_google" {
		t.Errorf("expected action_create_event_google, got %s", action)
	}
}

func TestDetectAction_InvalidJSON(t *testing.T) {
	logger.InitConsole()
	mock := &mockAI{
		Response: `Not a JSON at all`,
	}
	_, err := actions.DetectAction(mock, "Do something")
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
