package internalhttp

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	listen     string
	logger     Logger
	app        Application
	httpServer *http.Server
}

type Logger interface { // TODO
	Info(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	ListEvents(ctx context.Context, period storage.Period) ([]storage.Event, error)
}

// CreateEventRequest represents the request structure for creating an event
type CreateEventRequest struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Start       *time.Time `json:"start,omitempty"`
	End         *time.Time `json:"end,omitempty"`
	AllDay      float64    `json:"all_day,omitempty"`
	Clinic      *string    `json:"clinic,omitempty"`
	UserID      *int       `json:"user_id,omitempty"`
	Service     *string    `json:"service,omitempty"`
}

// CreateEventResponse represents the response structure for creating an event
type CreateEventResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ListEventsResponse represents the response structure for listing events
type ListEventsResponse struct {
	Success bool            `json:"success"`
	Events  []storage.Event `json:"events,omitempty"`
	Error   string          `json:"error,omitempty"`
}

func NewServer(logger Logger, app Application, listen string) *Server {
	return &Server{listen: listen, logger: logger, app: app}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy OK"))
	})

	// Add the create event endpoint
	mux.HandleFunc("/api/create", s.handleCreateEvent)
	// Add the list events endpoint
	mux.HandleFunc("/api/events", s.handleListEvents)

	handler := loggingMiddleware(s.logger)(mux)

	s.httpServer = &http.Server{
		Addr:              s.listen,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second, // for example
		IdleTimeout:       30 * time.Second, // Explicitly set
		WriteTimeout:      15 * time.Second, // Recommended
	}

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.httpServer.IdleTimeout)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(done)
	}()

	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	<-done
	return nil
}

// handleCreateEvent handles POST requests to create a new event
func (s *Server) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Create the event using the application layer
	event := storage.Event{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Start:       req.Start,
		End:         req.End,
		AllDay:      req.AllDay,
		Clinic:      req.Clinic,
		UserID:      req.UserID,
		Service:     req.Service,
	}
	err := s.app.CreateEvent(r.Context(), event)
	if err != nil {
		response := CreateEventResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CreateEventResponse{
		Success: true,
		Message: "Event created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// handleListEvents handles GET requests to list all events
func (s *Server) handleListEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, we'll list all events without period filtering
	// You can extend this to accept query parameters for period filtering
	events, err := s.app.ListEvents(r.Context(), storage.PeriodAll)
	if err != nil {
		response := ListEventsResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ListEventsResponse{
		Success: true,
		Events:  events,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	s.logger.Info("Shutting down HTTP server...")
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.logger.Info("Error during shutdown: " + err.Error())
	} else {
		s.logger.Info("HTTP server shutdown complete.")
	}
	return err
}

// TODO
