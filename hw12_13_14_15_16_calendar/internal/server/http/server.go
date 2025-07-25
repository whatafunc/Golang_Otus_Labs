package internalhttp

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	listen     string
	logger     Logger
	httpServer *http.Server
}

type Logger interface { // TODO
	Info(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application, listen string) *Server {
	_ = app // TODO
	return &Server{listen: listen, logger: logger}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy OK"))
	})

	handler := loggingMiddleware(s.logger)(mux)

	s.httpServer = &http.Server{
		Addr:              s.listen,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second, // for example
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
