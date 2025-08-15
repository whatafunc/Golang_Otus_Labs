package calendarGRPC

import (
	"context"
	"fmt"
	"time"

	calendarpb "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/calendarGRPC/pb"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/app"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func fromProtoEvent(pe *calendarpb.Event) storage.Event {
	if pe == nil {
		return storage.Event{}
	}
	return storage.Event{
		ID:          int(pe.Id),
		Title:       pe.Title,
		Description: pe.Description,
		Start:       parseTimePtr(pe.StartTime),
		End:         parseTimePtr(pe.EndTime),
		AllDay:      float64(pe.AllDay),
		Clinic: func() *string {
			if pe.Clinic != "" {
				return &pe.Clinic
			}
			return nil
		}(),
		UserID: func() *int {
			if pe.UserId != 0 {
				uid := int(pe.UserId)
				return &uid
			}
			return nil
		}(),
		Service: func() *string {
			if pe.Service != "" {
				return &pe.Service
			}
			return nil
		}(),
	}
}

func toProtoEvent(ev storage.Event) *calendarpb.Event {
	return &calendarpb.Event{
		Id:          int32(ev.ID),
		Title:       ev.Title,
		Description: ev.Description,
		StartTime:   formatTimePtr(ev.Start),
		EndTime:     formatTimePtr(ev.End),
		AllDay:      int32(ev.AllDay),
		Clinic: func() string {
			if ev.Clinic != nil {
				return *ev.Clinic
			}
			return ""
		}(),
		UserId: func() int32 {
			if ev.UserID != nil {
				return int32(*ev.UserID)
			}
			return 0
		}(),
		Service: func() string {
			if ev.Service != nil {
				return *ev.Service
			}
			return ""
		}(),
	}
}

func toProtoEvents(events []storage.Event) []*calendarpb.Event {
	result := make([]*calendarpb.Event, len(events))
	for i, ev := range events {
		result[i] = toProtoEvent(ev)
	}
	return result
}

// --- RPC Implementations ---

func (s *EventServer) CreateEvent(ctx context.Context, req *calendarpb.CreateEventRequest) (*calendarpb.CreateEventResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}
	if err := s.application.CreateEvent(ctx, fromProtoEvent(req.Event)); err != nil {
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
		return &calendarpb.GetEventResponse{
			Success: false,
			Event:   &calendarpb.Event{},
			Error:   fmt.Sprintf("event not found: %v", err),
		}, nil
	}

	return &calendarpb.GetEventResponse{
		Success: true,
		Event:   toProtoEvent(ev),
	}, nil
}

func (s *EventServer) ListEventsDay(ctx context.Context, req *calendarpb.ListEventsRequest) (*calendarpb.ListEventsResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}

	events, err := s.application.ListEvents(ctx, storage.PeriodDay)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events for the day: %v", err)
	}

	return &calendarpb.ListEventsResponse{
		Success: true,
		Events:  toProtoEvents(events),
	}, nil
}

func (s *EventServer) ListEventsWeek(ctx context.Context, req *calendarpb.ListEventsRequest) (*calendarpb.ListEventsResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}

	events, err := s.application.ListEvents(ctx, storage.PeriodWeek)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events for the week: %v", err)
	}

	return &calendarpb.ListEventsResponse{
		Success: true,
		Events:  toProtoEvents(events),
	}, nil
}

func (s *EventServer) ListEventsMonth(ctx context.Context, req *calendarpb.ListEventsRequest) (*calendarpb.ListEventsResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}

	events, err := s.application.ListEvents(ctx, storage.PeriodMonth)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events for the Month: %v", err)
	}

	return &calendarpb.ListEventsResponse{
		Success: true,
		Events:  toProtoEvents(events),
	}, nil
}

func (s *EventServer) UpdateEvent(ctx context.Context, req *calendarpb.UpdateEventRequest) (*calendarpb.UpdateEventResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}
	if req.Event == nil {
		return &calendarpb.UpdateEventResponse{
			Success: false,
			Error:   "no event data provided",
		}, nil
	}
	if err := s.application.UpdateEvent(ctx, fromProtoEvent(req.Event)); err != nil {
		return &calendarpb.UpdateEventResponse{
			Success: false,
			Error:   fmt.Sprintf("event not found: %v", err),
		}, nil
	}
	return &calendarpb.UpdateEventResponse{Success: true}, nil
}

func (s *EventServer) DeleteEvent(ctx context.Context, req *calendarpb.DeleteEventRequest) (*calendarpb.DeleteEventResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}
	if err := s.application.DeleteEvent(ctx, int(req.Id)); err != nil {
		return &calendarpb.DeleteEventResponse{
			Success: false,
			Error:   fmt.Sprintf("event not found: %v", err),
		}, nil
	}
	return &calendarpb.DeleteEventResponse{Success: true}, nil
}
