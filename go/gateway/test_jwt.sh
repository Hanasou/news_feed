#!/bin/bash

# JWT Authentication Test Script for GraphQL Server
# Make sure your server is running on localhost:8080

echo "=== JWT Authentication Test for GraphQL Server ==="
echo

# First, let's try to access a protected query without authentication
echo "1. Testing access without authentication..."
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"query": "{ todos { id text done } }"}' \
  http://localhost:8080/query
echo -e "\n"

# Authenticate and get tokens
echo "2. Authenticating user..."
AUTH_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { authenticateUser(input: { identifier: \"testuser\", password: \"testpass\" }) { accessToken refreshToken user { id name email role } } }"}' \
  http://localhost:8080/query)

echo "Auth Response: $AUTH_RESPONSE"
echo

# Extract access token (this is a simple extraction - in practice you'd use jq or similar)
ACCESS_TOKEN=$(echo $AUTH_RESPONSE | grep -o '"accessToken":"[^"]*' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
    echo "Failed to get access token. Make sure the server is running and the authentication mutation works."
    exit 1
fi

echo "Access Token: $ACCESS_TOKEN"
echo

# Now try to access the protected query with authentication
echo "3. Testing access with authentication..."
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"query": "{ todos { id text done userD } }"}' \
  http://localhost:8080/query
echo -e "\n"

# Test creating a todo with authentication
echo "4. Creating a todo with authentication..."
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"query": "mutation { createTodo(input: { text: \"Test todo\", userId: \"user123\" }) { id text done userD } }"}' \
  http://localhost:8080/query
echo -e "\n"

# Test accessing users (admin only)
echo "5. Testing admin-only access (should fail for regular user)..."
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"query": "{ users { id name email role } }"}' \
  http://localhost:8080/query
echo -e "\n"

echo "=== Test completed ==="
