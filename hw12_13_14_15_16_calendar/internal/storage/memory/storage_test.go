//nolint:depguard // allowed for tests
package memorystorage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/storage"
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
	events, err := store.ListEvents(context.Background(), storage.PeriodAll)
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

func TestStorage_ContextCancellation(t *testing.T) {
	s := New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Immediate cancellation

	err := s.CreateEvent(ctx, storage.Event{})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context cancel error, got %v", err)
	}
}

func TestStorage_Timeout(t *testing.T) {
	s := New()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()
	time.Sleep(200 * time.Microsecond) // Ensure timeout expires

	_, err := s.GetEvent(ctx, 1)
	if err == nil {
		t.Errorf("expected context error, got nil")
	} else if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context canceled or deadline exceeded error, got %v", err)
	}
}

func TestDeleteEvent(t *testing.T) {
	s := New()
	ctx := context.Background()

	// Create test event
	err := s.CreateEvent(ctx, storage.Event{Title: "Meeting"})
	require.NoError(t, err)

	// Delete existing event
	err = s.DeleteEvent(ctx, 1)
	require.NoError(t, err)

	// Verify deletion
	_, err = s.GetEvent(ctx, 1)
	require.ErrorIs(t, err, ErrNotFound)

	// Test deleting non-existent event
	err = s.DeleteEvent(ctx, 999)
	require.ErrorIs(t, err, ErrNotFound)

	// Test context cancellation
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	err = s.DeleteEvent(cancelCtx, 1)
	require.ErrorIs(t, err, context.Canceled)
}
