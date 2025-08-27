package calendargrpc

import (
	"context"
	"testing"

	calendarpb "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/calendarGRPC/pb"
)

func TestCreateEventReturnsErrorWhenAppIsNil(t *testing.T) {
	server := &EventServer{application: nil}
	req := &calendarpb.CreateEventRequest{Event: &calendarpb.Event{}}

	_, err := server.CreateEvent(context.Background(), req)
	if err == nil {
		t.Fatal("expected error when application is nil, got nil")
	}
}
