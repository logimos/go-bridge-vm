#!/bin/bash

# Interactive Intent Recognition Test Script
# This script demonstrates how the system asks for missing information

BASE_URL="http://localhost:8080/api/v1"

echo "=== Interactive Intent Recognition Test ==="
echo "Testing how the system handles missing required information"
echo

# Function to test intent extraction and show interactive behavior
test_interactive() {
    local test_name="$1"
    local text="$2"
    
    echo "Test: $test_name"
    echo "Input: \"$text\""
    
    response=$(curl -s -X POST "$BASE_URL/intent" \
        -H "Content-Type: application/json" \
        -d "{\"text\": \"$text\"}")
    
    echo "Response:"
    echo "$response" | jq .
    
    # Check if there are follow-up questions
    follow_up_count=$(echo "$response" | jq '.intent.follow_up | length')
    if [ "$follow_up_count" -gt 0 ]; then
        echo
        echo "🔍 Interactive Behavior Detected!"
        echo "Missing required fields:"
        echo "$response" | jq -r '.intent.missing[]' | sed 's/^/  - /'
        echo
        echo "Follow-up questions:"
        echo "$response" | jq -r '.intent.follow_up[]' | sed 's/^/  • /'
        echo
        echo "💡 The system is asking for missing information to complete the intent."
    else
        echo
        echo "✅ Complete - All required information provided!"
    fi
    
    echo
    echo "---"
    echo
}

# Check if server is running
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo "❌ Server is not running. Please start the server with enhanced local provider:"
    echo
    echo "export AI_PROVIDER=enhanced_local"
    echo "export INTENT_CONFIG_PATH=configs/personal_assistant.json"
    echo "go run main.go"
    exit 1
fi

echo "✅ Server is running. Testing Interactive Behavior..."
echo

# Test CreateEvent with missing information
echo "=== CreateEvent Tests (Missing Information) ==="
test_interactive "Create Event - Missing Title" "create calendar event"
test_interactive "Create Event - Missing Date/Time" "create calendar event for team meeting"
test_interactive "Create Event - Missing Time" "schedule appointment with doctor tomorrow"
test_interactive "Create Event - Complete" "create calendar event for team meeting tomorrow at 2pm"

# Test CreateContact with missing information
echo "=== CreateContact Tests (Missing Information) ==="
test_interactive "Create Contact - Missing Name" "create a new contact"
test_interactive "Create Contact - Missing Email" "add contact John Smith"
test_interactive "Create Contact - Complete" "create a new contact named Alice with email alice@example.com"

# Test CreateTask with missing information
echo "=== CreateTask Tests (Missing Information) ==="
test_interactive "Create Task - Missing Title" "create new task"
test_interactive "Create Task - Complete" "create new task called buy groceries"

# Test complex scenarios
echo "=== Complex Interactive Scenarios ==="
test_interactive "Partial Event Info" "I need to schedule something"
test_interactive "Partial Contact Info" "add someone to my contacts"
test_interactive "Partial Task Info" "I need to create a task"

echo "=== Interactive Test Complete ==="
echo
echo "🎯 Key Features Demonstrated:"
echo "  • Detects missing required fields"
echo "  • Generates contextual follow-up questions"
echo "  • Shows completion status"
echo "  • Provides confidence scores"
echo
echo "💡 This makes the system much more user-friendly and interactive!" 