package app

import (
	"context"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	GetEvent(ctx context.Context, id int) (storage.Event, error)
	ListEvents(ctx context.Context) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
