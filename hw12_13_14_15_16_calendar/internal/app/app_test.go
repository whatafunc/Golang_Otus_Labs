//nolint:depguard
package app

import (
	"context"
	"errors"
	"testing"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/logger"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrNotFound      = errors.New("event not found")
	ErrContextCancel = errors.New("operation canceled")
	ErrDuplicate     = errors.New("event with this ID already exists")
)

// fakeStorage is a complete mock of storageInterface for testing.
type fakeStorage struct {
	events map[int]storage.Event
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{
		events: make(map[int]storage.Event),
	}
}

func (f *fakeStorage) CreateEvent(ctx context.Context, event storage.Event) error {
	select {
	case <-ctx.Done():
		return ErrContextCancel
	default:
	}
	// Add this duplicate check ↓
	if _, exists := f.events[event.ID]; exists {
		return ErrDuplicate
	}
	f.events[event.ID] = event
	return nil
}

func (f *fakeStorage) GetEvent(ctx context.Context, id int) (storage.Event, error) {
	select {
	case <-ctx.Done():
		return storage.Event{}, ErrContextCancel
	default:
	}

	event, ok := f.events[id]
	if !ok {
		return storage.Event{}, ErrNotFound
	}
	return event, nil
}

func (f *fakeStorage) ListEvents(ctx context.Context, period storage.Period) ([]storage.Event, error) {
	_ = period
	select {
	case <-ctx.Done():
		return nil, ErrContextCancel
	default:
	}

	list := make([]storage.Event, 0, len(f.events))
	for _, e := range f.events {
		list = append(list, e)
	}
	return list, nil
}

func (f *fakeStorage) DeleteEvent(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ErrContextCancel
	default:
	}

	if _, exists := f.events[id]; !exists {
		return ErrNotFound
	}
	delete(f.events, id)
	return nil
}

func TestApp_CreateEvent(t *testing.T) {
	t.Parallel()

	// Setup
	log := logger.New("")
	fakeStore := newFakeStorage()
	app := &App{log: log, store: fakeStore}

	testCases := []struct {
		name    string
		id      int
		title   string
		wantErr bool
	}{
		{"success case", 1, "Test Event", false},
		{"duplicate ID", 1, "Duplicate", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			if tc.wantErr {
				// First create the event to force duplicate
				_ = fakeStore.CreateEvent(ctx, storage.Event{ID: tc.id, Title: "Existing"})
			}

			// Exercise
			err := app.CreateEvent(ctx, tc.id, tc.title)

			// Verify
			if tc.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("CreateEvent returned error: %v", err)
			}

			storedEvent, err := fakeStore.GetEvent(ctx, tc.id)
			if err != nil {
				t.Fatalf("GetEvent returned error: %v", err)
			}
			if storedEvent.ID != tc.id || storedEvent.Title != tc.title {
				t.Errorf("stored event mismatch: got %v, want ID=%d Title=%q", storedEvent, tc.id, tc.title)
			}
		})
	}
}

func TestApp_DeleteEvent(t *testing.T) {
	t.Parallel()

	// Setup
	log := logger.New("")
	fakeStore := newFakeStorage()
	app := &App{log: log, store: fakeStore}

	testCases := []struct {
		name        string
		preCreate   bool
		id          int
		expectError error
	}{
		{"success case", true, 1, nil},
		{"not found", false, 1, ErrNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			if tc.preCreate {
				_ = fakeStore.CreateEvent(ctx, storage.Event{ID: tc.id, Title: "Test Event"})
			}

			// Exercise
			err := app.store.DeleteEvent(ctx, tc.id)

			// Verify
			if tc.expectError != nil {
				if !errors.Is(err, tc.expectError) {
					t.Errorf("expected error %v, got %v", tc.expectError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("DeleteEvent failed: %v", err)
			}

			// Verify deletion
			_, err = fakeStore.GetEvent(ctx, tc.id)
			if !errors.Is(err, ErrNotFound) {
				t.Errorf("expected ErrNotFound after deletion, got: %v", err)
			}
		})
	}
}
