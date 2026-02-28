//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	_ = os.MkdirAll(filepath.Join(tempDir, ".context"), 0o750)

	// Point rc to this temp dir's .context
	rc.Reset()

	return tempDir, func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	}
}

func TestLoadWebhook_NoKey(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	url, err := LoadWebhook()
	if err != nil {
		t.Fatalf("LoadWebhook() error = %v", err)
	}
	if url != "" {
		t.Errorf("LoadWebhook() = %q, want empty", url)
	}
}

func TestLoadWebhook_NoFile(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create key but no encrypted file
	keyPath := filepath.Join(tempDir, ".context", config.FileContextKey)
	_ = os.WriteFile(keyPath, make([]byte, 32), 0o600)

	url, err := LoadWebhook()
	if err != nil {
		t.Fatalf("LoadWebhook() error = %v", err)
	}
	if url != "" {
		t.Errorf("LoadWebhook() = %q, want empty", url)
	}
}

func TestLoadWebhook_RoundTrip(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	want := "https://example.com/webhook?token=secret123"

	if err := SaveWebhook(want); err != nil {
		t.Fatalf("SaveWebhook() error = %v", err)
	}

	got, err := LoadWebhook()
	if err != nil {
		t.Fatalf("LoadWebhook() error = %v", err)
	}
	if got != want {
		t.Errorf("LoadWebhook() = %q, want %q", got, want)
	}
}

func TestEventAllowed_Nil(t *testing.T) {
	if EventAllowed("anything", nil) {
		t.Error("EventAllowed(anything, nil) = true, want false (opt-in only)")
	}
}

func TestEventAllowed_Empty(t *testing.T) {
	if EventAllowed("anything", []string{}) {
		t.Error("EventAllowed(anything, []) = true, want false (opt-in only)")
	}
}

func TestEventAllowed_Match(t *testing.T) {
	if !EventAllowed("loop", []string{"loop", "nudge"}) {
		t.Error("EventAllowed(loop, [loop nudge]) = false, want true")
	}
}

func TestEventAllowed_NoMatch(t *testing.T) {
	if EventAllowed("test", []string{"loop", "nudge"}) {
		t.Error("EventAllowed(test, [loop nudge]) = true, want false")
	}
}

func TestSend_NoWebhook(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	// No webhook configured — should noop without error
	err := Send("test", "hello", "session-1", nil)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestSend_EventFiltered(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Set up a server that should NOT be called
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		called = true
	}))
	defer ts.Close()

	// Configure webhook
	if err := SaveWebhook(ts.URL); err != nil {
		t.Fatalf("SaveWebhook() error = %v", err)
	}

	// Configure events filter to only allow "loop"
	rcContent := "notify:\n  events:\n    - loop\n"
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0o600)
	rc.Reset()

	// Send event "test" which is NOT in the allowed list
	err := Send("test", "hello", "session-1", nil)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if called {
		t.Error("server was called despite event being filtered")
	}
}

func TestSend_Payload(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	var received map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
	}))
	defer ts.Close()

	if err := SaveWebhook(ts.URL); err != nil {
		t.Fatalf("SaveWebhook() error = %v", err)
	}

	// Configure events to allow "loop"
	rcContent := "notify:\n  events:\n    - loop\n"
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0o600)
	rc.Reset()

	ref := NewTemplateRef("check-context-size", "window",
		map[string]any{"Percentage": 82, "TokenCount": "164k"})
	sendErr := Send("loop", "Loop completed after 5 iterations", "abc123", ref)
	if sendErr != nil {
		t.Fatalf("Send() error = %v", sendErr)
	}

	if received["event"] != "loop" {
		t.Errorf("Event = %v, want %q", received["event"], "loop")
	}
	if received["message"] != "Loop completed after 5 iterations" {
		t.Errorf("Message = %v, want %q", received["message"], "Loop completed after 5 iterations")
	}
	if received["session_id"] != "abc123" {
		t.Errorf("SessionID = %v, want %q", received["session_id"], "abc123")
	}
	if received["timestamp"] == nil || received["timestamp"] == "" {
		t.Error("Timestamp is empty")
	}
	if received["project"] == nil || received["project"] == "" {
		t.Error("Project is empty")
	}

	// Assert structured detail
	detail, ok := received["detail"].(map[string]any)
	if !ok {
		t.Fatalf("Detail is not an object: %T = %v", received["detail"], received["detail"])
	}
	if detail["hook"] != "check-context-size" {
		t.Errorf("Detail.hook = %v, want %q", detail["hook"], "check-context-size")
	}
	if detail["variant"] != "window" {
		t.Errorf("Detail.variant = %v, want %q", detail["variant"], "window")
	}
	vars, ok := detail["variables"].(map[string]any)
	if !ok {
		t.Fatalf("Detail.variables is not an object: %T", detail["variables"])
	}
	if vars["Percentage"] != float64(82) {
		t.Errorf("Detail.variables.Percentage = %v, want 82", vars["Percentage"])
	}
	if vars["TokenCount"] != "164k" {
		t.Errorf("Detail.variables.TokenCount = %v, want %q", vars["TokenCount"], "164k")
	}
}

func TestSend_NilDetailOmitted(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	var received map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
	}))
	defer ts.Close()

	if err := SaveWebhook(ts.URL); err != nil {
		t.Fatalf("SaveWebhook() error = %v", err)
	}

	rcContent := "notify:\n  events:\n    - test\n"
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0o600)
	rc.Reset()

	sendErr := Send("test", "hello", "session-1", nil)
	if sendErr != nil {
		t.Fatalf("Send() error = %v", sendErr)
	}

	if _, exists := received["detail"]; exists {
		t.Errorf("Detail should be omitted when nil, got: %v", received["detail"])
	}
}

func TestSend_HTTPErrorIgnored(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	if err := SaveWebhook(ts.URL); err != nil {
		t.Fatalf("SaveWebhook() error = %v", err)
	}

	// Configure events to allow "test"
	rcContent := "notify:\n  events:\n    - test\n"
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0o600)
	rc.Reset()

	// Should not return error even on HTTP 500
	err := Send("test", "hello", "session-1", nil)
	if err != nil {
		t.Fatalf("Send() error = %v, want nil (fire-and-forget)", err)
	}
}
