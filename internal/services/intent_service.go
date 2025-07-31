package services

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"strings"

	"myllm/internal/models"
)

// IntentService handles intent recognition logic
type IntentService struct {
	aiProvider AIProvider
	patterns   map[string]*regexp.Regexp
}

// NewIntentService creates a new intent service instance
func NewIntentService() *IntentService {
	// Create AI provider configuration
	config := AIProviderConfig{
		ProviderType: getEnv("AI_PROVIDER", "openai"),
		Model:        getEnv("AI_MODEL", ""),
		Temperature:  getFloatEnvVar("AI_TEMPERATURE", 0.1),
		MaxTokens:    getIntEnvVar("AI_MAX_TOKENS", 1000),
		BaseURL:      getEnv("AI_BASE_URL", ""),
		APIKey:       getEnv("OPENAI_API_KEY", ""),
	}

	// Create AI provider factory
	factory := NewAIProviderFactory(config)

	// Try to create the configured provider
	aiProvider, err := factory.CreateProvider()
	if err != nil {
		// Fallback to available providers
		availableProviders := factory.GetAvailableProviders()
		if len(availableProviders) > 0 {
			aiProvider = availableProviders[0]
		} else {
			// Last resort: create local provider
			aiProvider, _ = NewLocalAIProvider(config)
		}
	}

	// Initialize pattern matching for common intents
	patterns := map[string]*regexp.Regexp{
		"create_contact": regexp.MustCompile(`(?i)(create|add|new)\s+(?:contact|person)\s+(?:named\s+)?([a-zA-Z]+)(?:\s+with\s+email\s+([^\s]+))?`),
		"find_contact":   regexp.MustCompile(`(?i)(find|search|look\s+for)\s+(?:contact\s+)?([a-zA-Z]+)`),
		"update_contact": regexp.MustCompile(`(?i)(update|change|modify)\s+(?:contact\s+)?([a-zA-Z]+)`),
		"delete_contact": regexp.MustCompile(`(?i)(delete|remove|drop)\s+(?:contact\s+)?([a-zA-Z]+)`),
	}

	return &IntentService{
		aiProvider: aiProvider,
		patterns:   patterns,
	}
}

// ExtractIntent processes natural language and extracts structured intent
func (s *IntentService) ExtractIntent(ctx context.Context, text string) (*models.Intent, error) {
	normalizedText := models.NormalizeText(text)

	// Try pattern matching first for efficiency
	if intent := s.extractWithPatterns(normalizedText); intent != nil {
		return intent, nil
	}

	// Fall back to AI-based extraction
	return s.aiProvider.ExtractIntent(ctx, normalizedText)
}

// extractWithPatterns uses regex patterns to extract intent
func (s *IntentService) extractWithPatterns(text string) *models.Intent {
	for intentType, pattern := range s.patterns {
		matches := pattern.FindStringSubmatch(text)
		if len(matches) > 0 {
			return s.buildIntentFromMatches(intentType, matches)
		}
	}
	return nil
}

// buildIntentFromMatches constructs intent from regex matches
func (s *IntentService) buildIntentFromMatches(intentType string, matches []string) *models.Intent {
	intent := &models.Intent{
		Task: strings.ToUpper(intentType),
		Vars: make(map[string]interface{}),
	}

	switch intentType {
	case "create_contact":
		if len(matches) > 2 {
			intent.Vars["name"] = matches[2]
		}
		if len(matches) > 3 && matches[3] != "" {
			intent.Vars["email"] = matches[3]
		} else {
			intent.Vars["email"] = ""
		}
		intent.Vars["phone"] = ""

	case "find_contact", "update_contact", "delete_contact":
		if len(matches) > 2 {
			intent.Vars["name"] = matches[2]
		}
	}

	return intent
}

// GetAIProviderName returns the name of the current AI provider
func (s *IntentService) GetAIProviderName() string {
	if s.aiProvider != nil {
		return s.aiProvider.Name()
	}
	return "None"
}

// getEnvVar is a wrapper for os.Getenv to make testing easier
var getEnvVar = os.Getenv

// getIntEnvVar is a wrapper for getIntEnv to make testing easier
var getIntEnvVar = getIntEnv

// getFloatEnvVar is a wrapper for getFloatEnv to make testing easier
var getFloatEnvVar = getFloatEnv

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := getEnvVar(key); value != "" {
		return value
	}
	return fallback
}

// getIntEnv gets integer environment variable with fallback
func getIntEnv(key string, fallback int) int {
	if value := getEnvVar(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getFloatEnv gets float environment variable with fallback
func getFloatEnv(key string, fallback float64) float64 {
	if value := getEnvVar(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return fallback
}
