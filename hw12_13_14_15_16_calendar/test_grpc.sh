#!/bin/bash

# Test script for Calendar gRPC API
# Make sure the gRPC server is running on localhost:50051

SERVICE="calendar.CalendarService"
HOST="localhost:50051"

echo "=== Testing Calendar gRPC API ==="
echo

# 1. Create an event
echo "1. Creating an event..."
grpcurl -plaintext -d '{
  "event": {
    "title": "Test Meeting",
    "description": "This is a test event",
    "start_time": "2025-08-15T10:00:00Z",
    "end_time": "2025-08-15T11:00:00Z",
    "all_day": 0,
    "clinic": "Test Clinic",
    "user_id": 123,
    "service": "Consultation"
  }
}' $HOST $SERVICE/CreateEvent
echo

# 2. Creating malformed event (field names wrong)
echo "2. Creating malformed event..."
grpcurl -plaintext -d '{
  "event": {
    "title": "Malformed",
    "description": "This is a malformed event",
    "start_time": "WRONG_TIME_FORMAT",
    "end_time": "2025-08-01T11:00:00Z"
  }
}' $HOST $SERVICE/CreateEvent
echo

# 3. Create todays event
echo "3. Creating todays event..."
grpcurl -plaintext -d '{
  "event": {
    "title": "Another Test Event Today",
    "description": "todays test event",
    "start_time": "2025-08-13T14:00:00Z",
    "end_time": "2025-08-13T15:00:00Z"
  }
}' $HOST $SERVICE/CreateEvent
echo

# 4. List day events
echo "4. Listing day events..."
grpcurl -plaintext -d '{"date": "2025-08-15"}' $HOST $SERVICE/ListEventsDay
echo

# 5. List week events
echo "5. Listing week events..."
grpcurl -plaintext -d '{"date": "2025-08-10"}' $HOST $SERVICE/ListEventsWeek
echo

# 6. List month events
echo "6. Listing month events..."
grpcurl -plaintext -d '{"date": "2025-08-01"}' $HOST $SERVICE/ListEventsMonth
echo

# 7. Get non-existing event
echo "7. Getting non-existing event..."
grpcurl -plaintext -d '{"id": 999}' $HOST $SERVICE/GetEvent
echo

# 8. Get existing event (id=1)
echo "8. Getting existing event..."
grpcurl -plaintext -d '{"id": 1}' $HOST $SERVICE/GetEvent
echo

# 9. Delete non-existing event
echo "9. Deleting non-existing event..."
grpcurl -plaintext -d '{"id": 999}' $HOST $SERVICE/DeleteEvent
echo

# 10. Delete existing event (id=2)
echo "10. Deleting existing event..."
grpcurl -plaintext -d '{"id": 2}' $HOST $SERVICE/DeleteEvent
echo

# 11. Update event (non-existing)
echo "11. Updating non-existing event..."
grpcurl -plaintext -d '{
  "event": {
    "id": 999,
    "title": "Updated Title",
    "description": "Updated Description"
  }
}' $HOST $SERVICE/UpdateEvent
echo

# 12. Update event (existing id=1)
echo "12. Updating existing event..."
grpcurl -plaintext -d '{
  "event": {
    "id": 1,
    "title": "Updated Event Title",
    "description": "Updated description of the event",
    "start_time": "2025-08-15T09:00:00Z",
    "end_time": "2025-08-15T11:00:00Z",
    "all_day": 0,
    "clinic": "Main Clinic",
    "user_id": 123,
    "service": "Updated Service"
  }
}' $HOST $SERVICE/UpdateEvent
echo

echo "=== Test completed ==="
