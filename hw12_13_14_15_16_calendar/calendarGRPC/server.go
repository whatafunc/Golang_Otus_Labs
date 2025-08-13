package calendarGRPC

import (
	"context"
	"time"

	calendarpb "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/calendarGRPC/pb"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/app"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	GetEvent(ctx context.Context, id int) (storage.Event, error)
	ListEventsDay(ctx context.Context, date string) ([]storage.Event, error)
	ListEventsWeek(ctx context.Context, startDate string) ([]storage.Event, error)
	ListEventsMonth(ctx context.Context, startDate string) ([]storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id int) error
}

type EventServer struct {
	calendarpb.UnimplementedCalendarServiceServer
	application *app.App
}

func NewEventServer(application *app.App) *EventServer {
	return &EventServer{application: application}
}

// --- Helpers ---

func parseTimePtr(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func toProtoEvents(events []storage.Event) []*calendarpb.Event {
	protoEvents := make([]*calendarpb.Event, len(events))
	for i, ev := range events {
		protoEvents[i] = &calendarpb.Event{
			Id:          int32(ev.ID),
			Title:       ev.Title,
			Description: ev.Description,
			StartTime:   formatTimePtr(ev.Start),
			EndTime:     formatTimePtr(ev.End),
		}
	}
	return protoEvents
}

// --- RPC Implementations ---

func (s *EventServer) CreateEvent(ctx context.Context, req *calendarpb.CreateEventRequest) (*calendarpb.CreateEventResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}

	ev := storage.Event{
		ID:          int(req.Event.Id),
		Title:       req.Event.Title,
		Description: req.Event.Description,
		Start:       parseTimePtr(req.Event.StartTime),
		End:         parseTimePtr(req.Event.EndTime),
	}
	if err := s.application.CreateEvent(ctx, ev); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}
	return &calendarpb.CreateEventResponse{Success: true}, nil
}

func (s *EventServer) GetEvent(ctx context.Context, req *calendarpb.GetEventRequest) (*calendarpb.GetEventResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}

	ev, err := s.application.GetEvent(ctx, int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "event not found: %v", err)
	}
	return &calendarpb.GetEventResponse{
		Event: &calendarpb.Event{
			Id:          int32(ev.ID),
			Title:       ev.Title,
			Description: ev.Description,
			StartTime:   formatTimePtr(ev.Start),
			EndTime:     formatTimePtr(ev.End),
		},
	}, nil
}

// func (s *EventServer) ListEventsDay(ctx context.Context, req *calendarpb.ListEventsRequest) (*calendarpb.ListEventsResponse, error) {
// 	if s.application == nil {
// 		return nil, status.Error(codes.Internal, "application is not initialized")
// 	}

// 	events, err := s.application.ListEventsDay(ctx, req.Date)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to list day events: %v", err)
// 	}
// 	return &calendarpb.ListEventsResponse{Events: toProtoEvents(events)}, nil
// }

// func (s *EventServer) ListEventsWeek(ctx context.Context, req *calendarpb.ListEventsRequest) (*calendarpb.ListEventsResponse, error) {
// 	if s.application == nil {
// 		return nil, status.Error(codes.Internal, "application is not initialized")
// 	}

// 	events, err := s.application.ListEventsWeek(ctx, req.Date)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to list week events: %v", err)
// 	}
// 	return &calendarpb.ListEventsResponse{Events: toProtoEvents(events)}, nil
// }

// func (s *EventServer) ListEventsMonth(ctx context.Context, req *calendarpb.ListEventsRequest) (*calendarpb.ListEventsResponse, error) {
// 	if s.application == nil {
// 		return nil, status.Error(codes.Internal, "application is not initialized")
// 	}

// 	events, err := s.application.ListEventsMonth(ctx, req.Date)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to list month events: %v", err)
// 	}
// 	return &calendarpb.ListEventsResponse{Events: toProtoEvents(events)}, nil
// }

func (s *EventServer) UpdateEvent(ctx context.Context, req *calendarpb.UpdateEventRequest) (*calendarpb.UpdateEventResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}

	ev := storage.Event{
		ID:          int(req.Event.Id),
		Title:       req.Event.Title,
		Description: req.Event.Description,
		Start:       parseTimePtr(req.Event.StartTime),
		End:         parseTimePtr(req.Event.EndTime),
	}
	if err := s.application.UpdateEvent(ctx, ev); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update event: %v", err)
	}
	return &calendarpb.UpdateEventResponse{Success: true}, nil
}

func (s *EventServer) DeleteEvent(ctx context.Context, req *calendarpb.DeleteEventRequest) (*calendarpb.DeleteEventResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}

	if err := s.application.DeleteEvent(ctx, int(req.Id)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete event: %v", err)
	}
	return &calendarpb.DeleteEventResponse{Success: true}, nil
}
