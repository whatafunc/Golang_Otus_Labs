package memorystorage

import (
	"context"
	"sync"
)

type Storage struct {
	// TODO
	mu sync.RWMutex //nolint:unused
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(ctx context.Context, id, title string) error {
	// TODO: implement in-memory logic
	return nil
}

// TODO
