package logic

// Renderable defines an interface for commands that can be rendered as message
type Renderable interface {
	RenderMessage() string
}
