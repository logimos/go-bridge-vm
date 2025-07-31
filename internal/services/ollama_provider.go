package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"myllm/internal/models"
	"net/http"
	"time"
)

// OllamaProvider implements AIProvider for Ollama
type OllamaProvider struct {
	client *http.Client
	config AIProviderConfig
}

// OllamaRequest represents the request structure for Ollama API
type OllamaRequest struct {
	Model   string        `json:"model"`
	Prompt  string        `json:"prompt"`
	Stream  bool          `json:"stream"`
	Options OllamaOptions `json:"options,omitempty"`
}

// OllamaOptions represents Ollama generation options
type OllamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}

// OllamaResponse represents the response structure from Ollama API
type OllamaResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(config AIProviderConfig) (AIProvider, error) {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Test connection to Ollama
	testURL := baseURL + "/api/tags"
	resp, err := client.Get(testURL)
	if err != nil {
		return nil, fmt.Errorf("Ollama not available at %s: %w", baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama health check failed with status %d", resp.StatusCode)
	}

	return &OllamaProvider{
		client: client,
		config: config,
	}, nil
}

// ExtractIntent extracts intent using Ollama
func (p *OllamaProvider) ExtractIntent(ctx context.Context, text string) (*models.Intent, error) {
	baseURL := p.config.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	model := p.config.Model
	if model == "" {
		model = "llama2" // Default model
	}

	prompt := fmt.Sprintf(`Extract intent and variables from this text: "%s"

Return a JSON object with this structure:
{
  "task": "TASK_NAME",
  "vars": {
    "name": "extracted_name",
    "email": "extracted_email", 
    "phone": "extracted_phone"
  }
}

Common tasks: CREATE_CONTACT, FIND_CONTACT, UPDATE_CONTACT, DELETE_CONTACT
If no specific task is found, use "UNKNOWN" as task.
Extract any names, emails, or phone numbers you can find.

Respond with valid JSON only:`, text)

	request := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Options: OllamaOptions{
			Temperature: p.config.Temperature,
			NumPredict:  p.config.MaxTokens,
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Ollama request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/generate", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API error %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	// Parse AI response
	intent, err := models.FromJSON(ollamaResp.Response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Ollama response: %w", err)
	}

	return intent, nil
}

// Name returns the provider name
func (p *OllamaProvider) Name() string {
	return "Ollama"
}

// IsAvailable checks if Ollama is available
func (p *OllamaProvider) IsAvailable() bool {
	baseURL := p.config.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	resp, err := p.client.Get(baseURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
