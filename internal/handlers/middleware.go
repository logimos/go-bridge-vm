package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests with timing information
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request details
		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call next handler
		next.ServeHTTP(w, r)

		// Log response time
		duration := time.Since(start)
		log.Printf("Response: %s %s completed in %v", r.Method, r.URL.Path, duration)
	})
}

// HealthCheck handles health check requests
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "intent-recognition-api",
	}

	// Simple JSON response for health check
	json.NewEncoder(w).Encode(response)
}
