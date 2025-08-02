# Calendar API Endpoints

## Create Event

Creates a new event in the calendar.

**Endpoint:** `POST /api/create`

**Request Body:**
```json
{
  "id": 1,
  "title": "Meeting with Client",
  "description": "Discuss project requirements",
  "start": "2024-01-15T10:00:00Z",
  "end": "2024-01-15T11:00:00Z",
  "all_day": 0,
  "clinic": "Main Clinic",
  "user_id": 123,
  "service": "Consultation"
}
```

**Required Fields:**
- `title`: The title of the event (string)

**Optional Fields:**
- `id`: Event ID (integer)
- `description`: Event description (string)
- `start`: Start time (ISO 8601 format)
- `end`: End time (ISO 8601 format)
- `all_day`: Whether the event is all-day (float64)
- `clinic`: Associated clinic (string)
- `user_id`: Associated user ID (integer)
- `service`: Associated service (string)

**Response:**

Success (201 Created):
```json
{
  "success": true,
  "message": "Event created successfully"
}
```

Error (400 Bad Request):
```json
{
  "success": false,
  "error": "Title is required"
}
```

Error (500 Internal Server Error):
```json
{
  "success": false,
  "error": "Database error message"
}
```

**Example Usage:**

```bash
curl -X POST http://localhost:8081/api/create \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "title": "Team Meeting",
    "description": "Weekly team sync",
    "start": "2024-01-15T09:00:00Z",
    "end": "2024-01-15T10:00:00Z"
  }'
```

## List Events

Retrieves all events from the calendar.

**Endpoint:** `GET /api/events`

**Response:**

Success (200 OK):
```json
{
  "success": true,
  "events": [
    {
      "id": 1,
      "title": "Team Meeting",
      "description": "Weekly team sync",
      "start": "2024-01-15T09:00:00Z",
      "end": "2024-01-15T10:00:00Z",
      "all_day": 0,
      "clinic": null,
      "user_id": null,
      "service": null
    }
  ]
}
```

Error (500 Internal Server Error):
```json
{
  "success": false,
  "error": "Database error message"
}
```

**Example Usage:**

```bash
curl -X GET http://localhost:8080/api/events
```

## Health Check

**Endpoint:** `GET /health`

Returns the health status of the server.

**Response:**
```
Healthy OK
``` 