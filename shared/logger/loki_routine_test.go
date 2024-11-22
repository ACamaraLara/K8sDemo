package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
)

type MockLoggerOutput struct {
	LogQueue chan []byte
}

func (m *MockLoggerOutput) AddLog(log []byte) {
	m.LogQueue <- log
}

func TestStartLokiLogPublishRoutine(t *testing.T) {
	// Create a test HTTP server to simulate Loki service.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != lokiPostPath {
			t.Errorf("Expected request path %s, got %s", lokiPostPath, r.URL.Path)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Read and verify request body
		body := new(bytes.Buffer)
		body.ReadFrom(r.Body)

		var payload map[string]interface{}
		err := json.Unmarshal(body.Bytes(), &payload)
		if err != nil {
			t.Errorf("Invalid JSON body: %v", err)
		}
	}))
	defer server.Close()

	// Set the LOKI_URL environment variable to the test server URL
	os.Setenv("LOKI_URL", server.URL)
	defer os.Unsetenv("LOKI_URL")

	logger := InitServiceLogger(LoggerConfig{"INFO", 1000})

	if err := logger.StartLokiLogPublishRoutine(); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	// Send log to the queue managed my loki routine.
	log.Info().Msg("This is a test log")

	// Allow time for the goroutine to process the log
	time.Sleep(3 * time.Second)
}
