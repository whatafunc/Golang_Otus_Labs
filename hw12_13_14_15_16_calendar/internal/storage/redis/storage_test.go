//nolint:depguard
package redisstorage_test

import (
	"context"
	"errors"
	"testing"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/config"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"
	redisstorage "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage/redis"
)

func TestNew(t *testing.T) {
	cfg := config.RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	s := redisstorage.New(cfg)
	if s == nil {
		t.Fatal("New() returned nil storage")
	}
	// if s.Cfg() != cfg {
	// 	// We don't have Cfg() method but can check indirectly if needed
	// 	// Or skip this as cfg field is unexported
	// }
}

func TestCreateEvent(t *testing.T) {
	s := redisstorage.New(config.RedisConfig{})
	ctx := context.Background()
	event := storage.Event{ID: 1, Title: "test"}

	err := s.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("CreateEvent() returned unexpected error: %v", err)
	}
}

func TestGetEvent_NotImplemented(t *testing.T) {
	s := redisstorage.New(config.RedisConfig{})
	ctx := context.Background()

	_, err := s.GetEvent(ctx, 1)
	if err == nil {
		t.Fatal("GetEvent() expected error but got nil")
	}
	if err.Error() != "not implemented" {
		t.Fatalf("GetEvent() unexpected error: got %v, want %v", err, errors.New("not implemented"))
	}
}

func TestListEvents_NotImplemented(t *testing.T) {
	s := redisstorage.New(config.RedisConfig{})
	ctx := context.Background()

	_, err := s.ListEvents(ctx)
	if err == nil {
		t.Fatal("ListEvents() expected error but got nil")
	}
	if err.Error() != "not implemented" {
		t.Fatalf("ListEvents() unexpected error: got %v, want %v", err, errors.New("not implemented"))
	}
}
