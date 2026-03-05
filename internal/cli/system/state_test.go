//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestReadCounter_NoFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent")
	got := readCounter(path)
	if got != 0 {
		t.Errorf("readCounter(nonexistent) = %d, want 0", got)
	}
}

func TestReadCounter_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "counter")
	if writeErr := os.WriteFile(path, []byte("42"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	got := readCounter(path)
	if got != 42 {
		t.Errorf("readCounter(42) = %d, want 42", got)
	}
}

func TestReadCounter_InvalidContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "counter")
	if writeErr := os.WriteFile(path, []byte("abc"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	got := readCounter(path)
	if got != 0 {
		t.Errorf("readCounter(abc) = %d, want 0", got)
	}
}

func TestWriteCounter_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "counter")
	writeCounter(path, 17)
	got := readCounter(path)
	if got != 17 {
		t.Errorf("writeCounter/readCounter roundtrip = %d, want 17", got)
	}
}

func TestLogMessage_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")
	logMessage(logFile, "sess-1234", "hello world")

	if _, statErr := os.Stat(logFile); statErr != nil {
		t.Fatalf("logMessage did not create file: %v", statErr)
	}
	data, readErr := os.ReadFile(logFile)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(data), "hello world") {
		t.Errorf("log file does not contain message, got: %s", data)
	}
}

func TestLogMessage_Appends(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")
	logMessage(logFile, "sess-1234", "line one")
	logMessage(logFile, "sess-1234", "line two")

	data, readErr := os.ReadFile(logFile)
	if readErr != nil {
		t.Fatal(readErr)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d: %q", len(lines), string(data))
	}
}

func TestRotateLog_UnderLimit(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")
	if writeErr := os.WriteFile(logFile, []byte("small"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	rotateLog(logFile)

	backupFile := logFile + ".1"
	if _, statErr := os.Stat(backupFile); !os.IsNotExist(statErr) {
		t.Error("backup file should not exist for small log")
	}
}

func TestRotateLog_OverLimit(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")

	// Write content exceeding LogMaxBytes (1MB).
	bigContent := make([]byte, config.LogMaxBytes+1)
	for i := range bigContent {
		bigContent[i] = 'x'
	}
	if writeErr := os.WriteFile(logFile, bigContent, 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	rotateLog(logFile)

	backupFile := logFile + ".1"
	if _, statErr := os.Stat(backupFile); statErr != nil {
		t.Errorf("backup file should exist after rotation: %v", statErr)
	}
	// Original should have been renamed away.
	if _, statErr := os.Stat(logFile); !os.IsNotExist(statErr) {
		t.Error("original log file should not exist after rotation")
	}
}

func TestIsDailyThrottled_NoFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "marker")
	if isDailyThrottled(path) {
		t.Error("isDailyThrottled should return false when file does not exist")
	}
}

func TestIsDailyThrottled_Today(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "marker")
	touchFile(path)
	if !isDailyThrottled(path) {
		t.Error("isDailyThrottled should return true for file touched today")
	}
}

func TestIsDailyThrottled_Yesterday(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "marker")
	touchFile(path)

	yesterday := time.Now().Add(-25 * time.Hour)
	if chtErr := os.Chtimes(path, yesterday, yesterday); chtErr != nil {
		t.Fatal(chtErr)
	}
	if isDailyThrottled(path) {
		t.Error("isDailyThrottled should return false for file touched yesterday")
	}
}

func TestContextDirLine_Set(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CTX_DIR", dir)
	rc.Reset()

	got := contextDirLine()
	if !strings.HasPrefix(got, "Context: ") {
		t.Errorf("contextDirLine() = %q, want prefix %q", got, "Context: ")
	}
}

func TestTouchFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "touched")
	touchFile(path)
	if _, statErr := os.Stat(path); statErr != nil {
		t.Errorf("touchFile did not create file: %v", statErr)
	}
}

func TestWriteSessionStats_CreatesAndAppends(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	ctxDir := rc.ContextDir()
	if mkErr := os.MkdirAll(filepath.Join(ctxDir, config.DirState), 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}

	sid := "test-stats"

	// Write two entries
	writeSessionStats(sid, sessionStats{
		Timestamp: "2026-03-05T10:00:00Z", Prompt: 1,
		Tokens: 5000, Pct: 2, WindowSize: 200000,
		Model: "claude-opus-4-6", Event: "silent",
	})
	writeSessionStats(sid, sessionStats{
		Timestamp: "2026-03-05T10:01:00Z", Prompt: 2,
		Tokens: 12000, Pct: 6, WindowSize: 200000,
		Model: "claude-opus-4-6", Event: "silent",
	})

	statsPath := filepath.Join(ctxDir, config.DirState, "stats-"+sid+".jsonl")
	data, readErr := os.ReadFile(statsPath)
	if readErr != nil {
		t.Fatalf("stats file not created: %v", readErr)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	var first sessionStats
	if jsonErr := json.Unmarshal([]byte(lines[0]), &first); jsonErr != nil {
		t.Fatalf("failed to parse first line: %v", jsonErr)
	}
	if first.Prompt != 1 || first.Tokens != 5000 || first.Event != "silent" {
		t.Errorf("unexpected first entry: %+v", first)
	}

	var second sessionStats
	if jsonErr := json.Unmarshal([]byte(lines[1]), &second); jsonErr != nil {
		t.Fatalf("failed to parse second line: %v", jsonErr)
	}
	if second.Prompt != 2 || second.Tokens != 12000 {
		t.Errorf("unexpected second entry: %+v", second)
	}
}

func TestWriteSessionStats_OmitsEmptyModel(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	ctxDir := rc.ContextDir()
	if mkErr := os.MkdirAll(filepath.Join(ctxDir, config.DirState), 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}

	writeSessionStats("test-no-model", sessionStats{
		Timestamp: "2026-03-05T10:00:00Z", Prompt: 1,
		Event: "silent",
	})

	statsPath := filepath.Join(ctxDir, config.DirState, "stats-test-no-model.jsonl")
	data, readErr := os.ReadFile(statsPath)
	if readErr != nil {
		t.Fatalf("stats file not created: %v", readErr)
	}
	if strings.Contains(string(data), `"model"`) {
		t.Errorf("expected model field to be omitted when empty, got: %s", data)
	}
}
