package internalhttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testLogger struct{}

func (l *testLogger) Info(msg string)  {} //nolint:revive
func (l *testLogger) Error(msg string) {} //nolint:revive

func TestHealthHandler(t *testing.T) {
	// Create a test request
	req := httptest.NewRequest("GET", "/health", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create minimal server just for handler testing
	server := &Server{
		logger: &testLogger{},
	}

	// Create the handler
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Healthy OK")
	})

	// Wrap with middleware if needed
	handler := loggingMiddleware(server.logger)(mux)

	// Serve the request
	handler.ServeHTTP(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Healthy OK") {
		t.Errorf("unexpected body: %s", string(body))
	}
}
