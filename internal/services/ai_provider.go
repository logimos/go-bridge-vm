package services

import (
	"context"
	"myllm/internal/models"
)

// AIProvider defines the interface for different AI backends
type AIProvider interface {
	// ExtractIntent extracts structured intent from natural language text
	ExtractIntent(ctx context.Context, text string) (*models.Intent, error)

	// Name returns the provider name for logging/debugging
	Name() string

	// IsAvailable checks if the provider is available and configured
	IsAvailable() bool
}

// AIProviderConfig holds configuration for AI providers
type AIProviderConfig struct {
	ProviderType string  // "openai", "local", "ollama", etc.
	Model        string  // Model name
	Temperature  float64 // Temperature for generation
	MaxTokens    int     // Maximum tokens to generate
	BaseURL      string  // Base URL for API calls (for local providers)
	APIKey       string  // API key if required
}

// AIProviderFactory creates AI providers based on configuration
type AIProviderFactory struct {
	config AIProviderConfig
}

// NewAIProviderFactory creates a new factory with the given configuration
func NewAIProviderFactory(config AIProviderConfig) *AIProviderFactory {
	return &AIProviderFactory{
		config: config,
	}
}

// CreateProvider creates an AI provider based on the configuration
func (f *AIProviderFactory) CreateProvider() (AIProvider, error) {
	switch f.config.ProviderType {
	case "openai":
		return NewOpenAIProvider(f.config)
	case "ollama":
		return NewOllamaProvider(f.config)
	case "local":
		return NewLocalAIProvider(f.config)
	case "enhanced_local":
		configPath := getEnv("INTENT_CONFIG_PATH", "")
		return NewEnhancedLocalProvider(configPath)
	default:
		return NewOpenAIProvider(f.config) // Default fallback
	}
}

// GetAvailableProviders returns a list of available providers
func (f *AIProviderFactory) GetAvailableProviders() []AIProvider {
	var providers []AIProvider

	// Try OpenAI
	if openai, err := NewOpenAIProvider(f.config); err == nil && openai.IsAvailable() {
		providers = append(providers, openai)
	}

	// Try Ollama
	if ollama, err := NewOllamaProvider(f.config); err == nil && ollama.IsAvailable() {
		providers = append(providers, ollama)
	}

	// Try Enhanced Local AI
	configPath := getEnv("INTENT_CONFIG_PATH", "")
	if enhanced, err := NewEnhancedLocalProvider(configPath); err == nil && enhanced.IsAvailable() {
		providers = append(providers, enhanced)
	}

	// Try Local AI (fallback)
	if local, err := NewLocalAIProvider(f.config); err == nil && local.IsAvailable() {
		providers = append(providers, local)
	}

	return providers
}
