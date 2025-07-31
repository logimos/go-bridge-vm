package services

import (
	"context"
	"fmt"
	"myllm/internal/models"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements AIProvider for OpenAI
type OpenAIProvider struct {
	client *openai.Client
	config AIProviderConfig
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config AIProviderConfig) (AIProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	client := openai.NewClient(config.APIKey)

	return &OpenAIProvider{
		client: client,
		config: config,
	}, nil
}

// ExtractIntent extracts intent using OpenAI
func (p *OpenAIProvider) ExtractIntent(ctx context.Context, text string) (*models.Intent, error) {
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
Extract any names, emails, or phone numbers you can find.`, text)

	model := p.config.Model
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}

	resp, err := p.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       model,
			Temperature: float32(p.config.Temperature),
			MaxTokens:   p.config.MaxTokens,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an intent extraction assistant. Always respond with valid JSON only.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI extraction failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse AI response
	aiResponse := resp.Choices[0].Message.Content
	intent, err := models.FromJSON(aiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	return intent, nil
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "OpenAI"
}

// IsAvailable checks if OpenAI is available
func (p *OpenAIProvider) IsAvailable() bool {
	return p.config.APIKey != "" && p.client != nil
}
