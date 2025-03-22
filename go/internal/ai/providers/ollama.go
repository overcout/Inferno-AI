package providers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/overcout/Inferno-AI/internal/logger"
)

// OllamaEngine is an implementation of AIEngine using Ollama API.
type OllamaEngine struct {
	APIURL string
	Model  string
}

// NewOllamaEngine creates a new OllamaEngine instance.
func NewOllamaEngine(apiURL, model string) *OllamaEngine {
	return &OllamaEngine{
		APIURL: apiURL,
		Model:  model,
	}
}

// ProcessPrompt sends a prompt to Ollama and returns the raw text response.
func (o *OllamaEngine) ProcessPrompt(prompt string) (string, error) {
	logger.Info.Println("[Ollama] Sending prompt to model")

	reqBody := map[string]any{
		"model":  o.Model,
		"prompt": prompt,
		"stream": true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(o.APIURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "response") {
			var partial struct {
				Response string `json:"response"`
			}
			if err := json.Unmarshal([]byte(line), &partial); err == nil {
				fullResponse.WriteString(partial.Response)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.TrimSpace(fullResponse.String()), nil
}
