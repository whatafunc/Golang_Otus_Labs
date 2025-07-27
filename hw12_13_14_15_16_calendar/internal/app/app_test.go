//nolint:revive
package app

import (
	"context"
	"testing"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/logger"  //nolint:depguard
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

// fakeStorage is a simple in-memory mock of storageInterface for testing.
type fakeStorage struct {
	events map[int]storage.Event
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{
		events: make(map[int]storage.Event),
	}
}

func (f *fakeStorage) CreateEvent(ctx context.Context, event storage.Event) error {
	f.events[event.ID] = event
	return nil
}

func (f *fakeStorage) GetEvent(ctx context.Context, id int) (storage.Event, error) {
	event, ok := f.events[id]
	if !ok {
		return storage.Event{}, nil
	}
	return event, nil
}

func (f *fakeStorage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	list := make([]storage.Event, 0, 1)
	for _, e := range f.events {
		list = append(list, e)
	}
	return list, nil
}

func TestApp_CreateEvent(t *testing.T) {
	// Setup
	log := &logger.Logger{} // You can provide a real logger or a mock if needed
	fakeStore := newFakeStorage()

	app := &App{
		log:   log,
		store: fakeStore,
	}

	ctx := context.Background()
	testID := 1
	testTitle := "Test Event"

	// Exercise
	err := app.CreateEvent(ctx, testID, testTitle)
	if err != nil {
		t.Fatalf("CreateEvent returned error: %v", err)
	}

	// Verify
	storedEvent, err := fakeStore.GetEvent(ctx, testID)
	if err != nil {
		t.Fatalf("GetEvent returned error: %v", err)
	}
	if storedEvent.ID != testID || storedEvent.Title != testTitle {
		t.Errorf("stored event mismatch: got %v, want ID=%d Title=%q", storedEvent, testID, testTitle)
	}
}
