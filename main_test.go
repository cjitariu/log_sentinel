package main

import (
	"context"
	"log/slog"
	"log_sentinel/pkg/models"
	"testing"
	"time"
)

// MockLogger implements slog.Handler to capture logs for testing
type MockLogger struct {
	Warns  []string
	Errors []string
}

func (m *MockLogger) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (m *MockLogger) Handle(_ context.Context, r slog.Record) error {
	msg := r.Message
	switch r.Level {
	case slog.LevelWarn:
		m.Warns = append(m.Warns, msg)
	case slog.LevelError:
		m.Errors = append(m.Errors, msg)
	}
	return nil
}
func (m *MockLogger) WithAttrs(attrs []slog.Attr) slog.Handler { return m }
func (m *MockLogger) WithGroup(name string) slog.Handler       { return m }

func TestCheckJobDuration(t *testing.T) {
	mockHandler := &MockLogger{}
	logger := slog.New(mockHandler)

	tests := []struct {
		name     string
		duration time.Duration
		wantWarn bool
		wantErr  bool
	}{
		{"below", 4 * time.Minute, false, false},
		{"warning", 6 * time.Minute, true, false},
		{"error", 11 * time.Minute, false, true},
	}

	for _, tt := range tests {
		mockHandler.Warns = nil
		mockHandler.Errors = nil
		job := struct {
			Name     string
			Duration time.Duration
		}{tt.name, tt.duration}
		checkJobDuration(models.Job{Name: job.Name, Duration: job.Duration}, logger)

		if tt.wantWarn && len(mockHandler.Warns) == 0 {
			t.Errorf("%s: expected warning, got none", tt.name)
		}
		if !tt.wantWarn && len(mockHandler.Warns) > 0 {
			t.Errorf("%s: unexpected warning: %v", tt.name, mockHandler.Warns)
		}
		if tt.wantErr && len(mockHandler.Errors) == 0 {
			t.Errorf("%s: expected error, got none", tt.name)
		}
		if !tt.wantErr && len(mockHandler.Errors) > 0 {
			t.Errorf("%s: unexpected error: %v", tt.name, mockHandler.Errors)
		}
	}
}

func TestParseTimestamp(t *testing.T) {
	// Test time only format
	ts := "12:00:00"
	_, err := time.Parse("15:04:05", ts)
	if err != nil {
		t.Errorf("Failed to parse time-only timestamp: %v", err)
	}
}
