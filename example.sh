#!/bin/bash

# Intent Recognition API Examples
# Make sure the server is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1"

echo "=== Intent Recognition API Examples ==="
echo

# Health check
echo "1. Health Check:"
curl -s "$BASE_URL/health" | jq .
echo

# Create contact examples
echo "2. Create Contact - Simple:"
curl -s -X POST "$BASE_URL/intent" \
  -H "Content-Type: application/json" \
  -d '{"text": "create a new contact named bob"}' | jq .
echo

echo "3. Create Contact - With Email:"
curl -s -X POST "$BASE_URL/intent" \
  -H "Content-Type: application/json" \
  -d '{"text": "add contact named alice with email alice@example.com"}' | jq .
echo

echo "4. Create Contact - Complex:"
curl -s -X POST "$BASE_URL/intent" \
  -H "Content-Type: application/json" \
  -d '{"text": "I need to add John Smith to my contacts, his email is john.smith@company.com and phone is 555-1234"}' | jq .
echo

# Find contact examples
echo "5. Find Contact:"
curl -s -X POST "$BASE_URL/intent" \
  -H "Content-Type: application/json" \
  -d '{"text": "find contact sarah"}' | jq .
echo

# Update contact examples
echo "6. Update Contact:"
curl -s -X POST "$BASE_URL/intent" \
  -H "Content-Type: application/json" \
  -d '{"text": "update contact mike"}' | jq .
echo

# Delete contact examples
echo "7. Delete Contact:"
curl -s -X POST "$BASE_URL/intent" \
  -H "Content-Type: application/json" \
  -d '{"text": "delete contact jane"}' | jq .
echo

# Error handling example
echo "8. Error Handling - Empty Text:"
curl -s -X POST "$BASE_URL/intent" \
  -H "Content-Type: application/json" \
  -d '{"text": ""}' | jq .
echo

echo "=== Examples Complete ===" 