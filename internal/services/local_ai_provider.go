package services

import (
	"context"
	"myllm/internal/models"
	"regexp"
	"strings"
)

// LocalAIProvider implements AIProvider for local rule-based extraction
type LocalAIProvider struct {
	config AIProviderConfig
	// Enhanced patterns for more sophisticated local extraction
	enhancedPatterns map[string]*regexp.Regexp
	// Keywords for intent classification
	intentKeywords map[string][]string
}

// NewLocalAIProvider creates a new local AI provider
func NewLocalAIProvider(config AIProviderConfig) (AIProvider, error) {
	provider := &LocalAIProvider{
		config: config,
		enhancedPatterns: map[string]*regexp.Regexp{
			"email": regexp.MustCompile(`(?i)([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`),
			"phone": regexp.MustCompile(`(?i)(\+\d{1,3}[-.\s]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}`),
			"name":  regexp.MustCompile(`(?i)(?:named\s+|name\s+is\s+|call(?:ed)?\s+)([A-Z][a-z]+(?:\s+[A-Z][a-z]+)*)`),
		},
		intentKeywords: map[string][]string{
			"CREATE_CONTACT": {"create", "add", "new", "insert", "save", "store"},
			"FIND_CONTACT":   {"find", "search", "look", "locate", "get"},
			"UPDATE_CONTACT": {"update", "change", "modify", "edit", "alter"},
			"DELETE_CONTACT": {"delete", "remove", "drop", "erase", "clear"},
		},
	}

	return provider, nil
}

// ExtractIntent extracts intent using local rule-based processing
func (p *LocalAIProvider) ExtractIntent(ctx context.Context, text string) (*models.Intent, error) {
	normalizedText := strings.ToLower(strings.TrimSpace(text))

	// Determine intent based on keywords
	intent := p.classifyIntent(normalizedText)

	// Extract entities
	entities := p.extractEntities(text)

	// Build the intent structure
	result := &models.Intent{
		Task: intent,
		Vars: make(map[string]interface{}),
	}

	// Map extracted entities to variables
	if name, ok := entities["name"]; ok {
		result.Vars["name"] = name
	} else {
		result.Vars["name"] = ""
	}

	if email, ok := entities["email"]; ok {
		result.Vars["email"] = email
	} else {
		result.Vars["email"] = ""
	}

	if phone, ok := entities["phone"]; ok {
		result.Vars["phone"] = phone
	} else {
		result.Vars["phone"] = ""
	}

	return result, nil
}

// classifyIntent determines the intent based on keywords
func (p *LocalAIProvider) classifyIntent(text string) string {
	text = strings.ToLower(text)

	// Count keyword matches for each intent
	intentScores := make(map[string]int)

	for intent, keywords := range p.intentKeywords {
		score := 0
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				score++
			}
		}
		intentScores[intent] = score
	}

	// Find the intent with the highest score
	maxScore := 0
	bestIntent := "UNKNOWN"

	for intent, score := range intentScores {
		if score > maxScore {
			maxScore = score
			bestIntent = intent
		}
	}

	return bestIntent
}

// extractEntities extracts named entities from text
func (p *LocalAIProvider) extractEntities(text string) map[string]string {
	entities := make(map[string]string)

	// Extract email addresses
	if emailMatches := p.enhancedPatterns["email"].FindStringSubmatch(text); len(emailMatches) > 1 {
		entities["email"] = emailMatches[1]
	}

	// Extract phone numbers
	if phoneMatches := p.enhancedPatterns["phone"].FindStringSubmatch(text); len(phoneMatches) > 1 {
		entities["phone"] = phoneMatches[1]
	}

	// Extract names
	if nameMatches := p.enhancedPatterns["name"].FindStringSubmatch(text); len(nameMatches) > 1 {
		entities["name"] = nameMatches[1]
	} else {
		// Fallback: look for capitalized words that might be names
		words := strings.Fields(text)
		for _, word := range words {
			if len(word) > 1 && word[0] >= 'A' && word[0] <= 'Z' {
				// Check if it's not a common word
				if !p.isCommonWord(word) {
					entities["name"] = word
					break
				}
			}
		}
	}

	return entities
}

// isCommonWord checks if a word is a common word (not likely a name)
func (p *LocalAIProvider) isCommonWord(word string) bool {
	commonWords := map[string]bool{
		"the": true, "and": true, "or": true, "but": true, "in": true, "on": true, "at": true,
		"to": true, "for": true, "of": true, "with": true, "by": true, "from": true, "up": true,
		"about": true, "into": true, "through": true, "during": true, "before": true, "after": true,
		"above": true, "below": true, "between": true, "among": true, "within": true, "without": true,
		"contact": true, "person": true, "email": true, "phone": true, "name": true, "create": true,
		"add": true, "new": true, "find": true, "search": true, "update": true, "delete": true,
		"remove": true, "modify": true, "change": true, "edit": true, "save": true, "store": true,
	}

	return commonWords[strings.ToLower(word)]
}

// Name returns the provider name
func (p *LocalAIProvider) Name() string {
	return "Local AI (Rule-based)"
}

// IsAvailable checks if local AI is available (always true for local provider)
func (p *LocalAIProvider) IsAvailable() bool {
	return true // Local provider is always available
}
