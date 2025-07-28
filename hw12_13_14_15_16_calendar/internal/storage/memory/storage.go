//nolint:depguard // allowed - temporary
package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrNotFound      = errors.New("event not found")
	ErrContextCancel = errors.New("operation canceled")
)

type Storage struct {
	mu     sync.RWMutex
	events map[int]storage.Event
	nextID int
}

func New() *Storage {
	return &Storage{
		events: make(map[int]storage.Event),
		nextID: 1,
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	// Check context before acquiring lock
	select {
	case <-ctx.Done():
		return ErrContextCancel
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check context again after acquiring lock
	select {
	case <-ctx.Done():
		return ErrContextCancel
	default:
		// Simulate slow operation for demonstration
		time.Sleep(10 * time.Millisecond)

		event.ID = s.nextID
		s.nextID++
		s.events[event.ID] = event
		return nil
	}
}

func (s *Storage) GetEvent(ctx context.Context, id int) (storage.Event, error) {
	select {
	case <-ctx.Done():
		return storage.Event{}, ErrContextCancel
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return storage.Event{}, ErrContextCancel
	default:
		event, ok := s.events[id]
		if !ok {
			return storage.Event{}, ErrNotFound
		}
		return event, nil
	}
}

func (s *Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ErrContextCancel
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ErrContextCancel
	default:
		result := make([]storage.Event, 0, len(s.events))
		for _, event := range s.events {
			// Check context periodically during long operations
			select {
			case <-ctx.Done():
				return nil, ErrContextCancel
			default:
				result = append(result, event)
			}
		}
		return result, nil
	}
}
