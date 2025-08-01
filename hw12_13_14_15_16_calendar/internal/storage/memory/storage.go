//nolint:depguard // allowed - temporary
package memorystorage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"
)

var ErrNotFound = errors.New("event not found")

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
		return fmt.Errorf("context canceled after acquiring lock: %w", ctx.Err())
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check context again after acquiring lock
	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled after acquiring lock: %w", ctx.Err())
	default:
		// Simulate slow operation for demonstration
		// time.Sleep(10 * time.Millisecond)

		event.ID = s.nextID
		s.nextID++
		s.events[event.ID] = event
		return nil
	}
}

func (s *Storage) GetEvent(ctx context.Context, id int) (storage.Event, error) {
	select {
	case <-ctx.Done():
		return storage.Event{}, fmt.Errorf("context canceled after acquiring lock: %w", ctx.Err())
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return storage.Event{}, fmt.Errorf("context canceled after acquiring lock: %w", ctx.Err())
	default:
		event, ok := s.events[id]
		if !ok {
			return storage.Event{}, ErrNotFound
		}
		return event, nil
	}
}

// ListEvents returns events, optionally filtered by time period.
// PeriodAll returns all events, other periods filter relative to current time.
func (s *Storage) ListEvents(ctx context.Context, period storage.Period) ([]storage.Event, error) {
	// period := PeriodAll // default
	// if len(periods) > 0 {
	// 	period = periods[0]
	// }

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled after acquiring lock: %w", ctx.Err())
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context canceled after acquiring lock: %w", ctx.Err())
	default:
		now := time.Now() // Base period calculations on current time
		result := make([]storage.Event, 0, len(s.events))

		for _, event := range s.events {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("context canceled after acquiring lock: %w", ctx.Err())
			default:
				if period == storage.PeriodAll || matchesPeriod(*event.Start, now, period) {
					result = append(result, event)
				}
			}
		}
		return result, nil
	}
}

// matchesPeriod helper remains the same as previous implementation.
func matchesPeriod(eventTime, baseTime time.Time, period storage.Period) bool {
	switch period {
	case storage.PeriodAll:
		return true // Include all events when PeriodAll is specified.
	case storage.PeriodDay:
		return sameDay(eventTime, baseTime)
	case storage.PeriodWeek:
		return sameWeek(eventTime, baseTime)
	case storage.PeriodMonth:
		return sameMonth(eventTime, baseTime)
	default:
		return false
	}
}

// sameDay/sameWeek/sameMonth helpers.
func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func sameWeek(a, b time.Time) bool {
	a = a.AddDate(0, 0, -int(a.Weekday()))
	b = b.AddDate(0, 0, -int(b.Weekday()))
	return sameDay(a, b)
}

func sameMonth(a, b time.Time) bool {
	y1, m1, _ := a.Date()
	y2, m2, _ := b.Date()
	return y1 == y2 && m1 == m2
}

// DeleteEvent removes an event by ID. Returns ErrNotFound if event doesn't exist.
func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	// Check context before acquiring lock
	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled after acquiring lock: %w", ctx.Err())
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check context again after acquiring lock
	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled after acquiring lock: %w", ctx.Err())
	default:
		// Check if event exists
		if _, exists := s.events[id]; !exists {
			return ErrNotFound
		}

		// Simulate slow operation for demonstration
		// time.Sleep(10 * time.Millisecond)

		delete(s.events, id)
		return nil
	}
}

func (s *Storage) ClearAll(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled after acquiring lock: %w", ctx.Err())
	default:
		s.events = make(map[int]storage.Event)
		s.nextID = 1
		return nil
	}
}
