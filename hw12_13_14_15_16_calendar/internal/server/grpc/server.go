package calendargrpc

import (
	"context"
	"fmt"
	"time"

	calendarpb "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/calendarGRPC/pb"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/app"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventServer struct {
	calendarpb.UnimplementedCalendarServiceServer
	application *app.App
	logger      *logger.Logger
}

func NewEventServer(application *app.App, log *logger.Logger) *EventServer {
	return &EventServer{application: application, logger: log}
}

// NewGRPCServer creates a grpc.Server with logging interceptor and registers the EventServer.
func NewGRPCServer(app *app.App, log *logger.Logger) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(LoggingUnaryInterceptor(log)),
	}
	grpcServer := grpc.NewServer(opts...)

	eventServer := NewEventServer(app, log)
	calendarpb.RegisterCalendarServiceServer(grpcServer, eventServer)

	return grpcServer
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

func fromProtoEvent(pe *calendarpb.Event) (storage.Event, error) {
	if pe == nil {
		return storage.Event{}, fmt.Errorf("event is nil")
	}

	start := parseTimePtr(pe.Start)
	if start == nil {
		return storage.Event{}, fmt.Errorf("invalid start time: %q", pe.Start)
	}

	var end *time.Time
	if pe.End != "" {
		end = parseTimePtr(pe.End)
		if end == nil {
			return storage.Event{}, fmt.Errorf("invalid end time: %q", pe.End)
		}
	}

	return storage.Event{
			ID:          int(pe.Id),
			Title:       pe.Title,
			Description: pe.Description,
			Start:       start,
			End:         end,
			AllDay: func() float64 {
				if pe.AllDay {
					return 1
				}
				return 0
			}(),
			Clinic: &pe.Clinic,
			UserID: func() *int {
				if pe.UserId == 0 {
					return nil
				}
				uid := int(pe.UserId)
				return &uid
			}(),
			Service: &pe.Service,
		},
		nil
}

func toProtoEvent(ev storage.Event) *calendarpb.Event {
	return &calendarpb.Event{
		Id:          int32(ev.ID), //nolint:gosec
		Title:       ev.Title,
		Description: ev.Description,
		Start:       formatTimePtr(ev.Start),
		End:         formatTimePtr(ev.End),
		AllDay:      ev.AllDay != 0,
		Clinic: func() string {
			if ev.Clinic != nil {
				return *ev.Clinic
			}
			return ""
		}(),
		UserId: func() int32 {
			if ev.UserID != nil {
				return int32(*ev.UserID) //nolint:gosec
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

// Stub for period constants needed for filtering events.
const (
	PeriodDay = iota
	PeriodWeek
	PeriodMonth
)

// --- RPC Implementations ---

func (s *EventServer) CreateEvent(
	ctx context.Context,
	req *calendarpb.CreateEventRequest,
) (*calendarpb.CreateEventResponse, error) {
	// if s.application == nil { // moved to main.go - before server start step
	// 	s.logger.Error("application is not initialized")
	// 	return nil, status.Error(codes.Unavailable, "something went wrong, pls try again later")
	// }

	eventValidated, err := fromProtoEvent(req.Event)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to create event with bad data input for: %v", err))
		return nil, status.Errorf(codes.InvalidArgument, "something went wrong with recieved data")
	}

	if err := s.application.CreateEvent(ctx, eventValidated); err != nil {
		s.logger.Error(fmt.Sprintf("failed to create event: %v", err))
		return nil, status.Errorf(codes.Unavailable, "something went wrong, pls try again a bit later")
	}

	s.logger.Info("event created successfully")
	return &calendarpb.CreateEventResponse{Success: true}, nil
}

func (s *EventServer) GetEvent(ctx context.Context, req *calendarpb.GetEventRequest) (
	*calendarpb.GetEventResponse,
	error,
) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}

	ev, err := s.application.GetEvent(ctx, int(req.Id))
	if err != nil {
		return &calendarpb.GetEventResponse{
			Event: &calendarpb.Event{},
			Error: fmt.Sprintf("event not found: %v", err),
		}, nil
	}

	return &calendarpb.GetEventResponse{
		Event: toProtoEvent(ev),
	}, nil
}

func (s *EventServer) ListEventsDay(ctx context.Context, req *emptypb.Empty) (*calendarpb.ListEventsResponse, error) {
	_ = req
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}

	events, err := s.application.ListEvents(ctx, storage.PeriodDay)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events for the day: %v", err)
	}

	return &calendarpb.ListEventsResponse{
		Events: toProtoEvents(events),
	}, nil
}

func (s *EventServer) ListEventsWeek(ctx context.Context, req *emptypb.Empty) (*calendarpb.ListEventsResponse, error) {
	_ = req
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}

	events, err := s.application.ListEvents(ctx, storage.PeriodWeek)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events for the week: %v", err)
	}

	return &calendarpb.ListEventsResponse{
		Events: toProtoEvents(events),
	}, nil
}

func (s *EventServer) ListEventsMonth(ctx context.Context, req *emptypb.Empty) (*calendarpb.ListEventsResponse, error) {
	_ = req
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}

	events, err := s.application.ListEvents(ctx, storage.PeriodMonth)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events for the month: %v", err)
	}

	return &calendarpb.ListEventsResponse{
		Events: toProtoEvents(events),
	}, nil
}

func (s *EventServer) UpdateEvent(
	ctx context.Context,
	req *calendarpb.UpdateEventRequest,
) (*calendarpb.UpdateEventResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}
	if req.Event == nil {
		return &calendarpb.UpdateEventResponse{
			Success: false,
			Error:   "no event data provided",
		}, nil
	}

	eventValidated, err := fromProtoEvent(req.Event)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to create event with bad data input for: %v", err))
		return nil, status.Errorf(codes.InvalidArgument, "something went wrong with recieved data")
	}

	if err := s.application.UpdateEvent(ctx, eventValidated); err != nil {
		return &calendarpb.UpdateEventResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to update event: %v", err),
		}, nil
	}
	return &calendarpb.UpdateEventResponse{Success: true}, nil
}

func (s *EventServer) DeleteEvent(
	ctx context.Context,
	req *calendarpb.DeleteEventRequest,
) (*calendarpb.DeleteEventResponse, error) {
	if s.application == nil {
		return nil, status.Error(codes.Internal, "application is not initialized")
	}
	if err := s.application.DeleteEvent(ctx, int(req.Id)); err != nil {
		return &calendarpb.DeleteEventResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to delete event: %v", err),
		}, nil
	}
	return &calendarpb.DeleteEventResponse{Success: true}, nil
}
