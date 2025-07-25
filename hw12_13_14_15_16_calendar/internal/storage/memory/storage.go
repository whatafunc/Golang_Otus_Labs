//nolint:depguard // allowed - temporary
package memorystorage

import (
	"context"
	"errors"
	"sync"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[int]storage.Event
	nextID int
}

func New() *Storage {
	return &Storage{events: make(map[int]storage.Event), nextID: 1}
}

//nolint:revive // unused-parameter - temporary
func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	event.ID = s.nextID
	s.nextID++
	s.events[event.ID] = event
	return nil
}

//nolint:revive // unused-parameter - temporary
func (s *Storage) GetEvent(ctx context.Context, id int) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[id]
	if !ok {
		return storage.Event{}, ErrNotFound
	}
	return event, nil
}

//nolint:revive // unused-parameter - temporary
func (s *Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]storage.Event, 0, len(s.events))
	for _, event := range s.events {
		result = append(result, event)
	}
	return result, nil
}

var ErrNotFound = errors.New("event not found")
