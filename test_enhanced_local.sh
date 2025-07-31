#!/bin/bash

# Enhanced Local AI Provider Test Script
# This script tests the enhanced local AI provider with various intent examples

BASE_URL="http://localhost:8080/api/v1"

echo "=== Enhanced Local AI Provider Test ==="
echo

# Function to test intent extraction
test_intent() {
    local test_name="$1"
    local text="$2"
    
    echo "Test: $test_name"
    echo "Input: \"$text\""
    
    response=$(curl -s -X POST "$BASE_URL/intent" \
        -H "Content-Type: application/json" \
        -d "{\"text\": \"$text\"}")
    
    echo "Response:"
    echo "$response" | jq .
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

echo "✅ Server is running. Testing Enhanced Local AI Provider..."
echo

# Test Contact Management
echo "=== Contact Management Tests ==="
test_intent "Create Contact - Simple" "create a new contact named John Smith"
test_intent "Create Contact - With Email" "add contact Alice with email alice@example.com"
test_intent "Create Contact - Complex" "save contact info for Bob Johnson, his phone is 555-123-4567"
test_intent "Update Contact" "update contact details for Sarah Wilson"
test_intent "Delete Contact" "delete contact John Smith"
test_intent "Show Contact" "show contact details for Alice"

# Test Task Management
echo "=== Task Management Tests ==="
test_intent "Create Task" "create new task called buy groceries"
test_intent "Create Todo" "add todo item for meeting preparation"
test_intent "Today's Tasks" "what's today's tasks"
test_intent "Show Tasks" "show today's schedule"

# Test Event Management
echo "=== Event Management Tests ==="
test_intent "Create Event" "create calendar event for team meeting"
test_intent "Schedule Meeting" "schedule appointment with doctor"
test_intent "Add Event" "add meeting to calendar for tomorrow"

# Test Note Management
echo "=== Note Management Tests ==="
test_intent "Create Note" "create note about project ideas"
test_intent "Write Memo" "write memo for team meeting"
test_intent "Take Note" "take note of important points"

# Test Utility Functions
echo "=== Utility Function Tests ==="
test_intent "Weather Query" "what's the weather today"
test_intent "Weather Location" "what's the weather like in New York"
test_intent "Time Query" "what time is it"
test_intent "Calculator" "calculate 2 + 2"

# Test Edge Cases
echo "=== Edge Case Tests ==="
test_intent "Fuzzy Match" "add someone to my contacts"
test_intent "Synonym Test" "make a new task"
test_intent "Complex Query" "I need to create a contact for John Smith with email john@company.com and phone 555-987-6543"
test_intent "Unknown Intent" "random text that doesn't match any intent"

echo "=== Test Complete ==="
echo
echo "Provider Information:"
curl -s "$BASE_URL/health" | jq . 