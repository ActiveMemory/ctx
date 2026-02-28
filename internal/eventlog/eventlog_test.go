//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package eventlog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// setupTestDir creates a temporary directory, configures rc to use it,
// and returns a cleanup function.
func setupTestDir(t *testing.T, enableLog bool) string {
	t.Helper()
	dir := t.TempDir()

	rc.Reset()
	rc.OverrideContextDir(filepath.Join(dir, config.DirContext))

	// Write .ctxrc to control event_log.
	rcContent := "event_log: false\n"
	if enableLog {
		rcContent = "event_log: true\n"
	}
	if writeErr := os.WriteFile(
		filepath.Join(dir, config.FileContextRC), []byte(rcContent), config.PermFile,
	); writeErr != nil {
		t.Fatalf("failed to write .ctxrc: %v", writeErr)
	}

	// Change to temp dir so rc loads the .ctxrc.
	origDir, _ := os.Getwd()
	if chErr := os.Chdir(dir); chErr != nil {
		t.Fatalf("failed to chdir: %v", chErr)
	}
	rc.Reset() // force reload with new cwd

	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	return dir
}

func TestAppend_Disabled(t *testing.T) {
	dir := setupTestDir(t, false)
	logPath := filepath.Join(dir, config.DirContext, config.DirState, config.FileEventLog)

	Append("relay", "test message", "session-1", nil)

	if _, statErr := os.Stat(logPath); !os.IsNotExist(statErr) {
		t.Error("Append() created log file when event_log is disabled")
	}
}

func TestAppend_Basic(t *testing.T) {
	dir := setupTestDir(t, true)
	logPath := filepath.Join(dir, config.DirContext, config.DirState, config.FileEventLog)

	detail := notify.NewTemplateRef("qa-reminder", "gate", nil)
	Append("relay", "QA gate reminder", "session-1", detail)

	data, readErr := os.ReadFile(logPath) //nolint:gosec // test file
	if readErr != nil {
		t.Fatalf("failed to read log: %v", readErr)
	}

	var p notify.Payload
	if unmarshalErr := json.Unmarshal(data, &p); unmarshalErr != nil {
		t.Fatalf("failed to parse log line: %v", unmarshalErr)
	}

	if p.Event != "relay" {
		t.Errorf("Event = %q, want %q", p.Event, "relay")
	}
	if p.Message != "QA gate reminder" {
		t.Errorf("Message = %q, want %q", p.Message, "QA gate reminder")
	}
	if p.SessionID != "session-1" {
		t.Errorf("SessionID = %q, want %q", p.SessionID, "session-1")
	}
	if p.Detail == nil || p.Detail.Hook != "qa-reminder" {
		t.Errorf("Detail.Hook = %v, want %q", p.Detail, "qa-reminder")
	}
	if p.Timestamp == "" {
		t.Error("Timestamp is empty")
	}
}

func TestAppend_CreatesStateDir(t *testing.T) {
	dir := setupTestDir(t, true)
	stateDir := filepath.Join(dir, config.DirContext, config.DirState)

	// Verify state dir doesn't exist yet.
	if _, statErr := os.Stat(stateDir); !os.IsNotExist(statErr) {
		t.Fatal("state dir should not exist before Append")
	}

	Append("nudge", "test", "", nil)

	if _, statErr := os.Stat(stateDir); os.IsNotExist(statErr) {
		t.Error("Append() did not create state directory")
	}
}

func TestAppend_Rotation(t *testing.T) {
	dir := setupTestDir(t, true)
	logPath := filepath.Join(dir, config.DirContext, config.DirState, config.FileEventLog)
	prevPath := filepath.Join(dir, config.DirContext, config.DirState, config.FileEventLogPrev)

	// Create state dir and write a file that exceeds the max size.
	stateDir := filepath.Join(dir, config.DirContext, config.DirState)
	if mkErr := os.MkdirAll(stateDir, config.PermExec); mkErr != nil {
		t.Fatalf("failed to create state dir: %v", mkErr)
	}

	bigContent := strings.Repeat(`{"event":"relay","message":"filler"}`+"\n", 40000)
	if writeErr := os.WriteFile(logPath, []byte(bigContent), config.PermFile); writeErr != nil {
		t.Fatalf("failed to write big log: %v", writeErr)
	}

	// Append should trigger rotation.
	Append("relay", "after rotation", "", nil)

	// Previous file should exist with the big content.
	if _, statErr := os.Stat(prevPath); os.IsNotExist(statErr) {
		t.Error("rotation did not create events.1.jsonl")
	}

	// Current file should be small (just the new event).
	info, statErr := os.Stat(logPath)
	if statErr != nil {
		t.Fatalf("current log missing after rotation: %v", statErr)
	}
	if info.Size() > 1024 {
		t.Errorf("current log is %d bytes after rotation, expected small", info.Size())
	}
}

func TestAppend_RotationOverwrite(t *testing.T) {
	dir := setupTestDir(t, true)
	logPath := filepath.Join(dir, config.DirContext, config.DirState, config.FileEventLog)
	prevPath := filepath.Join(dir, config.DirContext, config.DirState, config.FileEventLogPrev)

	stateDir := filepath.Join(dir, config.DirContext, config.DirState)
	if mkErr := os.MkdirAll(stateDir, config.PermExec); mkErr != nil {
		t.Fatalf("failed to create state dir: %v", mkErr)
	}

	// Create an existing .1 file.
	if writeErr := os.WriteFile(prevPath, []byte("old rotated content\n"), config.PermFile); writeErr != nil {
		t.Fatalf("failed to write old .1 file: %v", writeErr)
	}

	// Write oversized current log.
	bigContent := strings.Repeat(`{"event":"relay","message":"filler"}`+"\n", 40000)
	if writeErr := os.WriteFile(logPath, []byte(bigContent), config.PermFile); writeErr != nil {
		t.Fatalf("failed to write big log: %v", writeErr)
	}

	Append("relay", "new event", "", nil)

	// The .1 file should now contain the rotated content, not "old rotated content".
	data, readErr := os.ReadFile(prevPath) //nolint:gosec // test file
	if readErr != nil {
		t.Fatalf("failed to read .1 file: %v", readErr)
	}
	if strings.Contains(string(data), "old rotated content") {
		t.Error("rotation did not overwrite existing .1 file")
	}
}

func TestQuery_NoFile(t *testing.T) {
	setupTestDir(t, true)

	events, queryErr := Query(QueryOpts{})
	if queryErr != nil {
		t.Fatalf("Query() error: %v", queryErr)
	}
	if len(events) != 0 {
		t.Errorf("Query() returned %d events, want 0", len(events))
	}
}

func TestQuery_FilterHook(t *testing.T) {
	setupTestDir(t, true)

	Append("relay", "qa gate", "s1", notify.NewTemplateRef("qa-reminder", "gate", nil))
	Append("relay", "context load", "s1", notify.NewTemplateRef("context-load-gate", "inject", nil))
	Append("nudge", "ceremonies", "s1", notify.NewTemplateRef("check-ceremonies", "both", nil))

	events, queryErr := Query(QueryOpts{Hook: "qa-reminder"})
	if queryErr != nil {
		t.Fatalf("Query() error: %v", queryErr)
	}
	if len(events) != 1 {
		t.Fatalf("Query(hook=qa-reminder) returned %d events, want 1", len(events))
	}
	if events[0].Message != "qa gate" {
		t.Errorf("Message = %q, want %q", events[0].Message, "qa gate")
	}
}

func TestQuery_FilterSession(t *testing.T) {
	setupTestDir(t, true)

	Append("relay", "session one", "s1", nil)
	Append("relay", "session two", "s2", nil)
	Append("relay", "session one again", "s1", nil)

	events, queryErr := Query(QueryOpts{Session: "s1"})
	if queryErr != nil {
		t.Fatalf("Query() error: %v", queryErr)
	}
	if len(events) != 2 {
		t.Errorf("Query(session=s1) returned %d events, want 2", len(events))
	}
}

func TestQuery_Last(t *testing.T) {
	setupTestDir(t, true)

	for i := 0; i < 20; i++ {
		Append("relay", "event", "", nil)
	}

	events, queryErr := Query(QueryOpts{Last: 5})
	if queryErr != nil {
		t.Fatalf("Query() error: %v", queryErr)
	}
	if len(events) != 5 {
		t.Errorf("Query(last=5) returned %d events, want 5", len(events))
	}
}

func TestQuery_IncludeRotated(t *testing.T) {
	dir := setupTestDir(t, true)
	stateDir := filepath.Join(dir, config.DirContext, config.DirState)
	if mkErr := os.MkdirAll(stateDir, config.PermExec); mkErr != nil {
		t.Fatalf("failed to create state dir: %v", mkErr)
	}

	// Write events to rotated file.
	prevPath := filepath.Join(stateDir, config.FileEventLogPrev)
	prevLine := `{"event":"relay","message":"old event","timestamp":"2026-01-01T00:00:00Z","project":"test"}` + "\n"
	if writeErr := os.WriteFile(prevPath, []byte(prevLine), config.PermFile); writeErr != nil {
		t.Fatalf("failed to write .1 file: %v", writeErr)
	}

	// Write event to current file.
	Append("relay", "new event", "", nil)

	// Without --all, only current events.
	events, _ := Query(QueryOpts{})
	if len(events) != 1 {
		t.Errorf("Query() without IncludeRotated returned %d events, want 1", len(events))
	}

	// With --all, both files.
	events, _ = Query(QueryOpts{IncludeRotated: true})
	if len(events) != 2 {
		t.Errorf("Query(IncludeRotated=true) returned %d events, want 2", len(events))
	}

	// Verify order: old event first, new event second.
	if events[0].Message != "old event" {
		t.Errorf("events[0].Message = %q, want %q", events[0].Message, "old event")
	}
}

func TestQuery_CorruptLine(t *testing.T) {
	dir := setupTestDir(t, true)
	stateDir := filepath.Join(dir, config.DirContext, config.DirState)
	if mkErr := os.MkdirAll(stateDir, config.PermExec); mkErr != nil {
		t.Fatalf("failed to create state dir: %v", mkErr)
	}

	logPath := filepath.Join(stateDir, config.FileEventLog)
	content := `{"event":"relay","message":"good","timestamp":"2026-01-01T00:00:00Z","project":"test"}
not valid json
{"event":"nudge","message":"also good","timestamp":"2026-01-02T00:00:00Z","project":"test"}
`
	if writeErr := os.WriteFile(logPath, []byte(content), config.PermFile); writeErr != nil {
		t.Fatalf("failed to write log: %v", writeErr)
	}

	events, queryErr := Query(QueryOpts{})
	if queryErr != nil {
		t.Fatalf("Query() error: %v", queryErr)
	}
	if len(events) != 2 {
		t.Errorf("Query() returned %d events, want 2 (corrupt line skipped)", len(events))
	}
}
