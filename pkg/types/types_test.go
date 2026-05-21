package model

import (
	"testing"

	"github.com/mycontroller-org/esphome_api/pkg/api"
)

func TestGetLogEntryConvertsCurrentSubscribeLogsResponse(t *testing.T) {
	entry, err := GetLogEntry(&api.SubscribeLogsResponse{
		Level:   api.LogLevel_LOG_LEVEL_WARN,
		Message: []byte("warning\n"),
	})
	if err != nil {
		t.Fatalf("GetLogEntry returned error: %v", err)
	}

	if entry.Level != LogLevelWarn {
		t.Fatalf("Level = %v, want %v", entry.Level, LogLevelWarn)
	}
	if entry.Message != "warning\n" {
		t.Fatalf("Message = %q, want %q", entry.Message, "warning\n")
	}
}
