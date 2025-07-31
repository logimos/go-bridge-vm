#!/bin/bash

# Intent Recognition API - Provider Demonstration
# This script demonstrates how to use different AI providers

BASE_URL="http://localhost:8080/api/v1"
TEST_TEXT="create a new contact named Alice with email alice@example.com"

echo "=== Intent Recognition API - Provider Demonstration ==="
echo

# Function to test intent extraction
test_intent() {
    local provider_name=$1
    echo "Testing with $provider_name..."
    
    response=$(curl -s -X POST "$BASE_URL/intent" \
        -H "Content-Type: application/json" \
        -d "{\"text\": \"$TEST_TEXT\"}")
    
    echo "Response: $response" | jq .
    echo
}

# Function to check health and provider info
check_health() {
    echo "Checking server health..."
    health_response=$(curl -s "$BASE_URL/health")
    echo "Health: $health_response" | jq .
    echo
}

# Check if server is running
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo "❌ Server is not running. Please start the server first:"
    echo "   go run main.go"
    echo
    echo "Or use one of these configurations:"
    echo
    echo "1. OpenAI (requires API key):"
    echo "   export AI_PROVIDER=openai"
    echo "   export OPENAI_API_KEY=your-key"
    echo "   go run main.go"
    echo
    echo "2. Ollama (requires Ollama installation):"
    echo "   export AI_PROVIDER=ollama"
    echo "   export AI_MODEL=llama2"
    echo "   go run main.go"
    echo
    echo "3. Local AI (no external dependencies):"
    echo "   export AI_PROVIDER=local"
    echo "   go run main.go"
    exit 1
fi

# Check health
check_health

# Test intent extraction
echo "Testing intent extraction with current provider..."
test_intent "current provider"

echo "=== Provider Configuration Examples ==="
echo

echo "1. OpenAI Configuration:"
echo "   export AI_PROVIDER=openai"
echo "   export OPENAI_API_KEY=your-openai-api-key"
echo "   export AI_MODEL=gpt-3.5-turbo"
echo "   export AI_TEMPERATURE=0.1"
echo

echo "2. Ollama Configuration:"
echo "   export AI_PROVIDER=ollama"
echo "   export AI_MODEL=llama2"
echo "   export AI_BASE_URL=http://localhost:11434"
echo "   export AI_TEMPERATURE=0.1"
echo

echo "3. Local AI Configuration:"
echo "   export AI_PROVIDER=local"
echo "   # No additional configuration needed"
echo

echo "=== Provider Comparison ==="
echo
echo "| Provider | Accuracy | Speed | Offline | Setup Complexity |"
echo "|----------|----------|-------|---------|------------------|"
echo "| OpenAI   | High     | Fast  | No      | Low (API key)    |"
echo "| Ollama   | Good     | Med   | Yes     | Medium (install) |"
echo "| Local    | Limited  | Fast  | Yes     | None             |"
echo

echo "=== Usage Examples ==="
echo

echo "Simple pattern matching (works with all providers):"
echo "curl -X POST $BASE_URL/intent \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"text\": \"create contact named bob\"}'"
echo

echo "Complex extraction (AI-powered):"
echo "curl -X POST $BASE_URL/intent \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"text\": \"I need to add John Smith to my contacts, his email is john.smith@company.com and phone is 555-1234\"}'"
echo

echo "=== Demo Complete ===" 