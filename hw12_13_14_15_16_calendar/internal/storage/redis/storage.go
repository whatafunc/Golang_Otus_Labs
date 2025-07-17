package redisstorage

import (
	"context"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/config"
)

type Storage struct {
	cfg config.RedisConfig
}

func New(cfg config.RedisConfig) *Storage {
	return &Storage{cfg: cfg}
}

// Implement the app.Storage interface methods as stubs
func (s *Storage) CreateEvent(ctx context.Context, id, title string) error {
	// TODO: implement Redis logic
	return nil
}
