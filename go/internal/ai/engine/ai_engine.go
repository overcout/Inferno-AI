package engine

// AIEngine defines a black-box interface for AI interaction.
type AIEngine interface {
	ProcessPrompt(prompt string) (string, error)
}
