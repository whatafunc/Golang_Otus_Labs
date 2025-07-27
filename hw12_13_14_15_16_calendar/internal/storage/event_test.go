package storage

import (
	"testing"
	"time"
)

func TestEventFields(t *testing.T) {
	// Setup example timestamps
	start := time.Now()
	end := start.Add(2 * time.Hour)

	userID := 123
	clinic := "Health Clinic"
	service := "Consultation"

	event := Event{
		ID:          1,
		Title:       "Medical Appointment",
		Description: "Annual check-up",
		Start:       &start,
		End:         &end,
		AllDay:      0,
		Clinic:      &clinic,
		UserID:      &userID,
		Service:     &service,
	}

	// Assertions
	if event.ID != 1 {
		t.Errorf("expected ID 1, got %d", event.ID)
	}

	if event.Title != "Medical Appointment" {
		t.Errorf("expected Title 'Medical Appointment', got %q", event.Title)
	}

	if event.Description != "Annual check-up" {
		t.Errorf("expected Description 'Annual check-up', got %q", event.Description)
	}

	if event.Start == nil || !event.Start.Equal(start) {
		t.Errorf("expected Start to be %v, got %v", start, event.Start)
	}

	if event.End == nil || !event.End.Equal(end) {
		t.Errorf("expected End to be %v, got %v", end, event.End)
	}

	if event.AllDay != 0 {
		t.Errorf("expected AllDay 0, got %v", event.AllDay)
	}

	if event.Clinic == nil || *event.Clinic != clinic {
		t.Errorf("expected Clinic %q, got %v", clinic, event.Clinic)
	}

	if event.UserID == nil || *event.UserID != userID {
		t.Errorf("expected UserID %d, got %v", userID, event.UserID)
	}

	if event.Service == nil || *event.Service != service {
		t.Errorf("expected Service %q, got %v", service, event.Service)
	}
}
