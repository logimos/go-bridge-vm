package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"myllm/internal/models"
	"myllm/internal/services"
)

// IntentHandler handles HTTP requests for intent extraction
type IntentHandler struct {
	intentService *services.IntentService
}

// NewIntentHandler creates a new intent handler
func NewIntentHandler(intentService *services.IntentService) *IntentHandler {
	return &IntentHandler{
		intentService: intentService,
	}
}

// ExtractIntent handles POST requests to extract intent from natural language
func (h *IntentHandler) ExtractIntent(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var request models.IntentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if request.Text == "" {
		respondWithError(w, http.StatusBadRequest, "Text field is required")
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Extract intent
	intent, err := h.intentService.ExtractIntent(ctx, request.Text)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to extract intent: "+err.Error())
		return
	}

	// Validate intent
	if err := intent.Validate(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid intent structure: "+err.Error())
		return
	}

	// Return success response
	response := models.IntentResponse{
		Success: true,
		Intent:  *intent,
	}

	respondWithJSON(w, http.StatusOK, response)
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	response := models.IntentResponse{
		Success: false,
		Error:   message,
	}
	respondWithJSON(w, statusCode, response)
}
