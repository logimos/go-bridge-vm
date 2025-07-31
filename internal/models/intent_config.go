package models

import (
	"encoding/json"
	"fmt"
	"os"
)

// IntentConfig represents a configurable intent recognition system
type IntentConfig struct {
	Domain     string                   `json:"domain"`     // e.g., "personal_assistant", "customer_support"
	Version    string                   `json:"version"`    // Config version
	Intents    map[string]IntentPattern `json:"intents"`    // Intent definitions
	Entities   map[string]EntityPattern `json:"entities"`   // Entity extraction patterns
	Synonyms   map[string][]string      `json:"synonyms"`   // Word synonyms for better matching
	Confidence map[string]float64       `json:"confidence"` // Confidence thresholds per intent
}

// IntentPattern defines how to recognize a specific intent
type IntentPattern struct {
	Description string   `json:"description"` // Human-readable description
	Keywords    []string `json:"keywords"`    // Primary keywords
	Phrases     []string `json:"phrases"`     // Common phrases
	Regex       []string `json:"regex"`       // Regex patterns
	Priority    int      `json:"priority"`    // Higher priority = more specific
	Variables   []string `json:"variables"`   // Expected variables to extract
	Required    []string `json:"required"`    // Required variables (will prompt if missing)
	Examples    []string `json:"examples"`    // Training examples
	FollowUp    []string `json:"follow_up"`   // Follow-up questions for missing info
}

// EntityPattern defines how to extract specific entities
type EntityPattern struct {
	Type        string   `json:"type"`        // Entity type (name, email, phone, etc.)
	Description string   `json:"description"` // Human-readable description
	Regex       []string `json:"regex"`       // Regex patterns for extraction
	Keywords    []string `json:"keywords"`    // Keywords that indicate this entity
	Examples    []string `json:"examples"`    // Example values
}

// LoadIntentConfig loads intent configuration from a JSON file
func LoadIntentConfig(filepath string) (*IntentConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config IntentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// Validate ensures the configuration is valid
func (c *IntentConfig) Validate() error {
	if c.Domain == "" {
		return fmt.Errorf("domain is required")
	}

	if len(c.Intents) == 0 {
		return fmt.Errorf("at least one intent must be defined")
	}

	// Validate each intent
	for intentName, intent := range c.Intents {
		if intent.Description == "" {
			return fmt.Errorf("intent %s: description is required", intentName)
		}
		if len(intent.Keywords) == 0 && len(intent.Phrases) == 0 && len(intent.Regex) == 0 {
			return fmt.Errorf("intent %s: must have at least keywords, phrases, or regex", intentName)
		}
	}

	return nil
}

// GetDefaultConfig returns a default configuration for personal assistant
func GetDefaultConfig() *IntentConfig {
	return &IntentConfig{
		Domain:  "personal_assistant",
		Version: "1.0.0",
		Intents: map[string]IntentPattern{
			"CREATE_CONTACT": {
				Description: "Create a new contact",
				Keywords:    []string{"create", "add", "new", "save"},
				Phrases:     []string{"create contact", "add contact", "new contact", "save contact"},
				Priority:    10,
				Variables:   []string{"name", "email", "phone"},
				Examples:    []string{"create a new contact named bob", "add contact alice with email alice@example.com"},
			},
			"FIND_CONTACT": {
				Description: "Find or search for a contact",
				Keywords:    []string{"find", "search", "look", "get"},
				Phrases:     []string{"find contact", "search contact", "look up contact"},
				Priority:    8,
				Variables:   []string{"name"},
				Examples:    []string{"find contact bob", "search for alice"},
			},
		},
		Entities: map[string]EntityPattern{
			"name": {
				Type:        "name",
				Description: "Person's name",
				Regex:       []string{`(?i)(?:named\s+|name\s+is\s+|call(?:ed)?\s+)([A-Z][a-z]+(?:\s+[A-Z][a-z]+)*)`},
				Keywords:    []string{"named", "name", "called"},
			},
			"email": {
				Type:        "email",
				Description: "Email address",
				Regex:       []string{`(?i)([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`},
				Keywords:    []string{"email", "e-mail", "mail"},
			},
			"phone": {
				Type:        "phone",
				Description: "Phone number",
				Regex:       []string{`(?i)(\+\d{1,3}[-.\s]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}`},
				Keywords:    []string{"phone", "telephone", "mobile", "cell"},
			},
		},
		Synonyms: map[string][]string{
			"create": {"add", "new", "save", "store", "insert"},
			"find":   {"search", "look", "locate", "get"},
			"update": {"change", "modify", "edit", "alter"},
			"delete": {"remove", "drop", "erase", "clear"},
		},
		Confidence: map[string]float64{
			"CREATE_CONTACT": 0.7,
			"FIND_CONTACT":   0.6,
			"UPDATE_CONTACT": 0.6,
			"DELETE_CONTACT": 0.6,
		},
	}
}
