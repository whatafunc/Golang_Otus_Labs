#!/bin/bash

# Test script for Calendar API endpoints
# Make sure the server is running on localhost

echo "=== Testing Calendar API ==="
# echo
# echo "CHECK health endpoint..."
# curl -X GET http://myapp.local/health

echo
# Test 1: Create an event
echo "1. Creating an event..."
curl -X POST http://myapp.local/api/create \
  -H "Content-Type: application/json" \
  -d '{
     "event": {
       "title": "Test chicken 3",
       "description": "This is a test event",
       "start": "2026-01-02T08:00:00Z",
       "end": "2026-01-02T11:50:00Z",
       "allDay": false,
       "clinic": "Test Clinic",
       "userId": 123,
       "service": "Consultation"
     }
  }' | jq '.'


# # Test !!!: Create a mailformed event
# echo "!. Creating mailformed event..."
# curl -X POST http://myapp.local/api/create \
#   -H "Content-Type: application/json" \
#   -d '{
#     "id": 1,
#     "title": "mailformed",
#     "description": "This is a mailformed event",
#     "start": "2025-01-15T10:00:00Z",
#     "end": "2025-01-15T11:00:00Z",
#     "all_day": 0,
#     "clinic": "Test mailformed Clinic",
#     "user_id": 123,
#     "service": "mailformed Consultation"
#   }' | jq '.'


# echo
# echo "2. Creating todays event..."
# curl -X POST http://myapp.local/api/create \
#   -H "Content-Type: application/json" \
#   -d '{
#     "id": 2,
#     "title": "Another Test Event todays",
#     "description": "todays test event",
#     "start": "2025-08-10T14:00:00Z",
#     "end": "2025-08-10T15:00:00Z"
#   }' | jq '.'

# echo
# echo "3. Listing all events..."
# curl -X GET http://myapp.local/api/events | jq '.'

#echo
#echo "3. Listing day events..."
#curl -X GET http://myapp.local/api/events?period=day | jq '.'

echo
echo "3.1. Listing day events..."
curl -X GET http://myapp.local/api/eventsDay | jq '.'

echo
echo "3.2. Listing Week events..."
curl -X GET http://myapp.local/api/eventsWeek | jq '.'

echo
echo "3.3. Listing Month events..."
curl -X GET http://myapp.local/api/eventsMonth | jq '.'

echo
echo "4.1 Testing get NonExisting event endpoint..."
curl -X GET http://myapp.local/api/get/26 | jq '.'

echo
echo "4.2 Testing get Existing event endpoint..."
curl -X GET http://myapp.local/api/get/2 | jq '.'

echo
echo "5. Testing delete endpoint..."
curl -X DELETE http://myapp.local/api/delete/1

echo
echo "6. Update event"
curl -X PUT http://myapp.local/api/update/2 \
  -H "Content-Type: application/json" \
  -d '{
      "event": {
        "title": "Test Meeting1-vera1",
        "description": "This is a test event1",
        "start_time": "2025-08-17T10:00:00Z",
        "end_time": "2025-08-17T11:00:00Z",
        "all_day": 0,
        "clinic": "Test Clinic12",
        "user_id": 123,
        "service": "Consultation1"
      }
  }'

# echo
# echo "7. Listing events after deletion..."
# curl -X GET http://myapp.local/api/events | jq '.'



echo
echo "=== Test completed ===" 