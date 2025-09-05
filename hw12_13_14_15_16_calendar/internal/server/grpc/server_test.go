package calendargrpc

import (
	"context"
	"testing"

	calendarpb "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/calendarGRPC/pb"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/logger"
)

func TestCreateEventReturnsErrorWhenAppIsNil(t *testing.T) {
	// Create a real logger instead of nil
	log := logger.New("info")
	server := &EventServer{application: nil, logger: log}
	req := &calendarpb.CreateEventRequest{Event: &calendarpb.Event{}}

	_, err := server.CreateEvent(context.Background(), req)
	if err == nil {
		t.Fatal("expected error when application is nil, got nil")
	}
}
