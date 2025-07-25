//nolint:revive // till next stage
package redisstorage

import (
	"context"
	"errors"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type Storage struct {
	cfg config.RedisConfig
}

func New(cfg config.RedisConfig) *Storage {
	return &Storage{cfg: cfg}
}

// Implement the app.Storage interface methods as stubs.
func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	// TODO: implement Redis logic
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id int) (storage.Event, error) {
	return storage.Event{}, errors.New("not implemented")
}

func (s *Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	return nil, errors.New("not implemented")
}
