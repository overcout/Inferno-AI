package logic

// Renderable defines an interface for commands that can be rendered as message
type Renderable interface {
	RenderMessage() string
}

// RenderText wraps a plain text string into a Renderable interface
type textResponse struct {
	Message string
}

func (t *textResponse) RenderMessage() string {
	return t.Message
}

// RenderText creates a Renderable object from a plain string
func RenderText(msg string) Renderable {
	return &textResponse{Message: msg}
}
