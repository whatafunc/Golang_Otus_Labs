//nolint:depguard // not finished app
package app

import (
	"context"
	"fmt"
	"os"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/config"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/logger"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage/memory"
	postgresstorage "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage/sql"
)

// App is the main application structure.
type App struct {
	log   *logger.Logger
	store storageInterface
}

// NewWithConfig creates and returns a new App instance based on the config.
func NewWithConfig(cfg config.Config, log *logger.Logger) *App {
	var store storageInterface

	switch cfg.Storage.Type {
	case "memory":
		store = memorystorage.New()
	// case "redis":
	// 	store = redisstorage.New(cfg.Storage.Redis)
	case "postgres":
		store = postgresstorage.New(cfg.Storage.Postgres)
	default:
		log.Error(fmt.Sprintf("unknown storage type: %s", cfg.Storage.Type))
		os.Exit(1)
	}

	return &App{
		log:   log,
		store: store,
	}
}

// storageInterface defines the expected behavior for all storage backends.
type storageInterface interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	GetEvent(ctx context.Context, id int) (storage.Event, error)
	ListEvents(ctx context.Context, period storage.Period) ([]storage.Event, error) // Updated
	DeleteEvent(ctx context.Context, id int) error
}

// CreateEvent adds a new event using the configured storage.
func (a *App) CreateEvent(ctx context.Context, id int, title string) error {
	event := storage.Event{
		ID:    id,
		Title: title,
	}
	return a.store.CreateEvent(ctx, event)
}

// Optionally, you could add other methods like GetEvent, ListEvents etc.
