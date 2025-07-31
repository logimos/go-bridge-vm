package services

import (
	"context"
	"testing"

	"myllm/internal/models"
)

func TestIntentService_ExtractIntent_PatternMatching(t *testing.T) {
	// Mock environment variables for testing
	originalGetEnv := getEnvVar
	originalGetIntEnv := getIntEnvVar
	originalGetFloatEnv := getFloatEnvVar

	getEnvVar = func(key string) string {
		switch key {
		case "AI_PROVIDER":
			return "local"
		case "AI_MODEL":
			return ""
		case "AI_BASE_URL":
			return ""
		case "OPENAI_API_KEY":
			return "test-key"
		default:
			return ""
		}
	}

	getIntEnvVar = func(key string, fallback int) int {
		return fallback
	}

	getFloatEnvVar = func(key string, fallback float64) float64 {
		return fallback
	}

	defer func() {
		getEnvVar = originalGetEnv
		getIntEnvVar = originalGetIntEnv
		getFloatEnvVar = originalGetFloatEnv
	}()

	service := NewIntentService()

	tests := []struct {
		name     string
		input    string
		expected *models.Intent
	}{
		{
			name:  "create contact with name",
			input: "create a new contact named bob",
			expected: &models.Intent{
				Task: "CREATE_CONTACT",
				Vars: map[string]interface{}{
					"name":  "bob",
					"email": "",
					"phone": "",
				},
			},
		},
		{
			name:  "create contact with name and email",
			input: "add contact named alice with email alice@example.com",
			expected: &models.Intent{
				Task: "CREATE_CONTACT",
				Vars: map[string]interface{}{
					"name":  "alice",
					"email": "alice@example.com",
					"phone": "",
				},
			},
		},
		{
			name:  "find contact",
			input: "find contact john",
			expected: &models.Intent{
				Task: "FIND_CONTACT",
				Vars: map[string]interface{}{
					"name": "john",
				},
			},
		},
		{
			name:  "update contact",
			input: "update contact sarah",
			expected: &models.Intent{
				Task: "UPDATE_CONTACT",
				Vars: map[string]interface{}{
					"name": "sarah",
				},
			},
		},
		{
			name:  "delete contact",
			input: "delete contact mike",
			expected: &models.Intent{
				Task: "DELETE_CONTACT",
				Vars: map[string]interface{}{
					"name": "mike",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := service.ExtractIntent(ctx, tt.input)

			if err != nil {
				t.Errorf("ExtractIntent() error = %v", err)
				return
			}

			if result.Task != tt.expected.Task {
				t.Errorf("Task = %v, want %v", result.Task, tt.expected.Task)
			}

			// Check variables
			for key, expectedValue := range tt.expected.Vars {
				if actualValue, exists := result.Vars[key]; !exists {
					t.Errorf("Missing variable %s", key)
				} else if actualValue != expectedValue {
					t.Errorf("Variable %s = %v, want %v", key, actualValue, expectedValue)
				}
			}
		})
	}
}

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "  Create   Contact   Named   Bob  ",
			expected: "create contact named bob",
		},
		{
			input:    "ADD CONTACT WITH EMAIL test@example.com",
			expected: "add contact with email test@example.com",
		},
		{
			input:    "find\tcontact\njohn",
			expected: "find contact john",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := models.NormalizeText(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeText() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIntent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		intent  *models.Intent
		wantErr bool
	}{
		{
			name: "valid intent",
			intent: &models.Intent{
				Task: "CREATE_CONTACT",
				Vars: map[string]interface{}{
					"name": "bob",
				},
			},
			wantErr: false,
		},
		{
			name: "missing task",
			intent: &models.Intent{
				Task: "",
				Vars: map[string]interface{}{
					"name": "bob",
				},
			},
			wantErr: true,
		},
		{
			name: "nil vars",
			intent: &models.Intent{
				Task: "CREATE_CONTACT",
				Vars: nil,
			},
			wantErr: false, // Should initialize empty map
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.intent.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Intent.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
