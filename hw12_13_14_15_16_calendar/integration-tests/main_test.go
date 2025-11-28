package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

// getBaseURL returns the base URL for the API, defaulting if the env var is not set.
func getBaseURL(t *testing.T) string {
	t.Helper()
	baseURL := os.Getenv("CALENDAR_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8888"
		t.Log("CALENDAR_API_URL not set, defaulting to", baseURL)
	}
	return baseURL
}

func TestCalendarAPI(t *testing.T) {
	baseURL := getBaseURL(t)
	t.Logf("Running tests against: %s", baseURL)

	// --- Test Health Check ---
	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/health")
		if err != nil {
			t.Fatalf("Health check request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 OK, got %d", resp.StatusCode)
		}

		var healthResponse struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&healthResponse); err != nil {
			t.Fatalf("Failed to decode health check response: %v", err)
		}

		if healthResponse.Status != "OK" {
			t.Errorf("Expected status 'OK', got '%s'", healthResponse.Status)
		}
	})

	// --- Test Create Event ---
	t.Run("CreateEvent", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)
		endTime := startTime.Add(1 * time.Hour)

		// This payload matches the structure of the CreateEventRequest protobuf message
		payload := map[string]interface{}{
			"event": map[string]interface{}{
				"title":       "Important Meeting",
				"description": "Discuss quarterly results.",
				"start":       startTime.Format(time.RFC3339),
				"end":         endTime.Format(time.RFC3339),
			},
		}

		eventJSON, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal event payload: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, baseURL+"/api/create", bytes.NewReader(eventJSON))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Create event request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
		}
	})
}
