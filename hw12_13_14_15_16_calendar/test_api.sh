#!/bin/bash

# Test script for Calendar API endpoints
# Make sure the server is running on localhost:8081

echo "=== Testing Calendar API ==="
echo

# Test 1: Create an event
echo "1. Creating an event..."
curl -X POST http://localhost:8081/api/create \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "title": "Test Meeting",
    "description": "This is a test event",
    "start": "2024-01-15T10:00:00Z",
    "end": "2024-01-15T11:00:00Z",
    "all_day": 0,
    "clinic": "Test Clinic",
    "user_id": 123,
    "service": "Consultation"
  }' | jq '.'

echo
echo "2. Creating another event..."
curl -X POST http://localhost:8081/api/create \
  -H "Content-Type: application/json" \
  -d '{
    "id": 2,
    "title": "Another Test Event",
    "description": "Second test event",
    "start": "2024-01-16T14:00:00Z",
    "end": "2024-01-16T15:00:00Z"
  }' | jq '.'

echo
echo "3. Listing all events..."
curl -X GET http://localhost:8081/api/events | jq '.'

echo
echo "4. Testing health endpoint..."
curl -X GET http://localhost:8081/health

echo
echo "=== Test completed ===" 