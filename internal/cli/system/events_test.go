//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// writeTestEvents writes test events to the event log file.
func writeTestEvents(t *testing.T, dir string, events []notify.Payload) {
	t.Helper()
	stateDir := filepath.Join(dir, config.DirState)
	if mkErr := os.MkdirAll(stateDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	logPath := filepath.Join(stateDir, config.FileEventLog)
	var buf bytes.Buffer
	for _, e := range events {
		line, _ := json.Marshal(e)
		buf.Write(line)
		buf.WriteByte('\n')
	}
	if writeErr := os.WriteFile(logPath, buf.Bytes(), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
}

func TestEventsCmd_Default(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CTX_DIR", dir)
	rc.Reset()

	ts := time.Now().UTC().Format(time.RFC3339)
	events := []notify.Payload{
		{Event: "relay", Message: "qa-reminder: QA gate emitted", Detail: notify.NewTemplateRef("qa-reminder", "gate", nil), Timestamp: ts, Project: "test"},
		{Event: "nudge", Message: "check-persistence: No context updated", Detail: notify.NewTemplateRef("check-persistence", "nudge", nil), Timestamp: ts, Project: "test"},
	}
	writeTestEvents(t, dir, events)

	cmd := eventsCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("events command failed: %v", runErr)
	}

	output := out.String()
	if !strings.Contains(output, "qa-reminder") {
		t.Errorf("expected qa-reminder in output, got: %s", output)
	}
	if !strings.Contains(output, "check-persistence") {
		t.Errorf("expected check-persistence in output, got: %s", output)
	}
}

func TestEventsCmd_JSON(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CTX_DIR", dir)
	rc.Reset()

	ts := time.Now().UTC().Format(time.RFC3339)
	events := []notify.Payload{
		{Event: "relay", Message: "test message", Timestamp: ts, Project: "test"},
	}
	writeTestEvents(t, dir, events)

	cmd := eventsCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--json"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("events --json failed: %v", runErr)
	}

	// Verify output is valid JSON
	var parsed notify.Payload
	if unmarshalErr := json.Unmarshal([]byte(strings.TrimSpace(out.String())), &parsed); unmarshalErr != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", unmarshalErr, out.String())
	}
	if parsed.Message != "test message" {
		t.Errorf("expected message 'test message', got %q", parsed.Message)
	}
}

func TestEventsCmd_NoLog(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CTX_DIR", dir)
	rc.Reset()

	cmd := eventsCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("events command failed: %v", runErr)
	}

	if !strings.Contains(out.String(), "No events logged.") {
		t.Errorf("expected 'No events logged.' message, got: %s", out.String())
	}
}

func TestEventsCmd_Filters(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CTX_DIR", dir)
	rc.Reset()

	ts := time.Now().UTC().Format(time.RFC3339)
	events := []notify.Payload{
		{Event: "relay", Message: "qa-reminder: gate", Detail: notify.NewTemplateRef("qa-reminder", "gate", nil), SessionID: "sess-1", Timestamp: ts, Project: "test"},
		{Event: "nudge", Message: "check-persistence: nudge", Detail: notify.NewTemplateRef("check-persistence", "nudge", nil), SessionID: "sess-2", Timestamp: ts, Project: "test"},
		{Event: "relay", Message: "qa-reminder: gate again", Detail: notify.NewTemplateRef("qa-reminder", "gate", nil), SessionID: "sess-2", Timestamp: ts, Project: "test"},
	}
	writeTestEvents(t, dir, events)

	// Filter by hook + session (intersection)
	cmd := eventsCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--hook", "qa-reminder", "--session", "sess-2"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("events with filters failed: %v", runErr)
	}

	output := out.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 filtered result, got %d: %s", len(lines), output)
	}
	if !strings.Contains(output, "qa-reminder") {
		t.Errorf("expected qa-reminder in output, got: %s", output)
	}
}

func TestFormatEventTimestamp(t *testing.T) {
	got := formatEventTimestamp("2026-02-27T22:39:31Z")
	if got == "" || got == "2026-02-27T22:39:31Z" {
		// Should be converted to local time (not the same as input unless UTC)
		// Just verify it doesn't return empty
		if got == "" {
			t.Error("formatEventTimestamp returned empty string")
		}
	}
}

func TestExtractHookName(t *testing.T) {
	tests := []struct {
		name    string
		payload notify.Payload
		want    string
	}{
		{
			name:    "from detail",
			payload: notify.Payload{Detail: notify.NewTemplateRef("qa-reminder", "gate", nil), Message: "qa-reminder: gate"},
			want:    "qa-reminder",
		},
		{
			name:    "from message prefix",
			payload: notify.Payload{Message: "check-persistence: nudge emitted"},
			want:    "check-persistence",
		},
		{
			name:    "no hook info",
			payload: notify.Payload{Message: "no colon here"},
			want:    "-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractHookName(tt.payload)
			if got != tt.want {
				t.Errorf("extractHookName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTruncateMessage(t *testing.T) {
	short := "short"
	if got := truncateMessage(short, 60); got != short {
		t.Errorf("truncateMessage(%q) = %q, want %q", short, got, short)
	}

	long := strings.Repeat("x", 100)
	got := truncateMessage(long, 60)
	if len(got) != 60 {
		t.Errorf("truncateMessage length = %d, want 60", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Error("truncated message should end with ...")
	}
}
