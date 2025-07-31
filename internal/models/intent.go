package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Intent represents the extracted intent and variables from natural language
type Intent struct {
	Task       string                 `json:"task"`
	Vars       map[string]interface{} `json:"vars"`
	Confidence float64                `json:"confidence,omitempty"`
	Missing    []string               `json:"missing,omitempty"`     // Required fields that are missing
	FollowUp   []string               `json:"follow_up,omitempty"`   // Questions to ask for missing info
	IsComplete bool                   `json:"is_complete,omitempty"` // Whether all required fields are present
}

// IntentRequest represents the incoming request to extract intent
type IntentRequest struct {
	Text string `json:"text" validate:"required"`
}

// IntentResponse represents the response with extracted intent
type IntentResponse struct {
	Success bool   `json:"success"`
	Intent  Intent `json:"intent,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ContactIntent represents a specific contact-related intent
type ContactIntent struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// Validate ensures the intent has required fields
func (i *Intent) Validate() error {
	if i.Task == "" {
		return fmt.Errorf("task is required")
	}
	if i.Vars == nil {
		i.Vars = make(map[string]interface{})
	}
	return nil
}

// ToJSON converts the intent to JSON string
func (i *Intent) ToJSON() (string, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return "", fmt.Errorf("failed to marshal intent: %w", err)
	}
	return string(data), nil
}

// FromJSON creates an intent from JSON string
func FromJSON(data string) (*Intent, error) {
	var intent Intent
	if err := json.Unmarshal([]byte(data), &intent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal intent: %w", err)
	}
	return &intent, nil
}

// NormalizeText cleans and normalizes input text for better processing
func NormalizeText(text string) string {
	// Convert to lowercase and trim whitespace
	normalized := strings.ToLower(strings.TrimSpace(text))

	// Remove extra whitespace
	normalized = strings.Join(strings.Fields(normalized), " ")

	return normalized
}
