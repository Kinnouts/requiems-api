#!/bin/bash

# Disposable Email API Testing Script
# Make sure the API server is running: go run apps/api/main.go

BASE_URL="http://localhost:8080/v1/email"

echo "=== Disposable Email API Tests ==="
echo ""

# Test 1: Check Single Email (Disposable)
echo "1. Testing single disposable email..."
curl -X POST "$BASE_URL/disposable/check" \
  -H "Content-Type: application/json" \
  -d '{"email": "user@tempmail.com"}' \
  -w "\n\n"

# Test 2: Check Single Email (Not Disposable)
echo "2. Testing single legitimate email..."
curl -X POST "$BASE_URL/disposable/check" \
  -H "Content-Type: application/json" \
  -d '{"email": "user@gmail.com"}' \
  -w "\n\n"

# Test 3: Batch Check
echo "3. Testing batch email check..."
curl -X POST "$BASE_URL/disposable/check-batch" \
  -H "Content-Type: application/json" \
  -d '{
    "emails": [
      "user1@tempmail.com",
      "user2@gmail.com",
      "user3@10minutemail.com",
      "user4@outlook.com"
    ]
  }' \
  -w "\n\n"

# Test 4: Check Domain
echo "4. Testing domain check..."
curl "$BASE_URL/disposable/domain/tempmail.com" \
  -w "\n\n"

# Test 5: Get Stats
echo "5. Testing stats endpoint..."
curl "$BASE_URL/disposable/stats" \
  -w "\n\n"

# Test 6: List Domains (first page)
echo "6. Testing domains list (first 10)..."
curl "$BASE_URL/disposable/domains?page=1&per_page=10" \
  -w "\n\n"

# Test 7: Error handling - missing email
echo "7. Testing error handling (missing email)..."
curl -X POST "$BASE_URL/disposable/check" \
  -H "Content-Type: application/json" \
  -d '{}' \
  -w "\n\n"

# Test 8: Error handling - batch limit exceeded
echo "8. Testing error handling (batch too large)..."
# Create array of 101 emails
emails=$(printf '"%s@test.com",' {1..101} | sed 's/,$//')
curl -X POST "$BASE_URL/disposable/check-batch" \
  -H "Content-Type: application/json" \
  -d "{\"emails\": [$emails]}" \
  -w "\n\n"

echo "=== Tests Complete ==="
