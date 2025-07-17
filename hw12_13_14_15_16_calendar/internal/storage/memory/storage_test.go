package memorystorage

import (
	"context"
	"testing"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"
)

func TestCreateAndGetEvent(t *testing.T) {
	store := New()
	event := storage.Event{
		Title:       "Test Event",
		Description: "A test event",
		AllDay:      1,
	}
	err := store.CreateEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	// The event should have ID 1
	got, err := store.GetEvent(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetEvent failed: %v", err)
	}
	if got.Title != event.Title || got.Description != event.Description || got.AllDay != event.AllDay {
		t.Errorf("GetEvent returned wrong data: got %+v, want %+v", got, event)
	}
}

func TestListEvents(t *testing.T) {
	store := New()
	for i := 0; i < 3; i++ {
		event := storage.Event{Title: "Event", Description: "Desc", AllDay: float64(i)}
		if err := store.CreateEvent(context.Background(), event); err != nil {
			t.Fatalf("CreateEvent failed: %v", err)
		}
	}
	events, err := store.ListEvents(context.Background())
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}
	if len(events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(events))
	}
}

func TestGetEvent_NotFound(t *testing.T) {
	store := New()
	_, err := store.GetEvent(context.Background(), 42)
	if err == nil {
		t.Errorf("Expected error for missing event, got nil")
	}
}
