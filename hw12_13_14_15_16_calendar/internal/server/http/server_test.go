package internalhttp

import (
	"context"
	"net/http"
	"testing"
	"time"
)

type testLogger struct {
	infos  []string
	errors []string
}

func (l *testLogger) Info(msg string)  { l.infos = append(l.infos, msg) }
func (l *testLogger) Error(msg string) { l.errors = append(l.errors, msg) }

type testApp struct{}

func TestServer_HealthEndpoint(t *testing.T) {
	logger := &testLogger{}
	app := &testApp{}
	server := NewServer(logger, app, ":8085") // Use a fixed port for test

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Start(ctx)

	// Wait a moment for server to start
	time.Sleep(200 * time.Millisecond)

	resp, err := http.Get("http://localhost:8085/health")
	if err != nil {
		t.Fatalf("failed to GET /health: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}

	// Shutdown server
	cancel()
	time.Sleep(100 * time.Millisecond)
}
