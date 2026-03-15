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

	"github.com/ActiveMemory/ctx/internal/config/crypto"
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
	keyPath := filepath.Join(tempDir, ".context", crypto.ContextKey)
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

	// No webhook configured: should noop without error
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

func TestSaveWebhook_Roundtrip(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	want := "https://hooks.example.com/notify?key=abc123"

	if saveErr := SaveWebhook(want); saveErr != nil {
		t.Fatalf("SaveWebhook() error = %v", saveErr)
	}

	got, loadErr := LoadWebhook()
	if loadErr != nil {
		t.Fatalf("LoadWebhook() error = %v", loadErr)
	}
	if got != want {
		t.Errorf("LoadWebhook() = %q, want %q", got, want)
	}
}

func TestLoadWebhook_CorruptedFile(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Save a valid webhook first (to create the key file).
	if saveErr := SaveWebhook("https://example.com"); saveErr != nil {
		t.Fatalf("SaveWebhook() error = %v", saveErr)
	}

	// Corrupt the encrypted file with garbage bytes.
	encPath := filepath.Join(tempDir, ".context", crypto.NotifyEnc)
	if writeErr := os.WriteFile(encPath, []byte("corrupted-garbage-data"), 0o600); writeErr != nil {
		t.Fatalf("WriteFile() error = %v", writeErr)
	}

	// Should return an error, not panic.
	_, loadErr := LoadWebhook()
	if loadErr == nil {
		t.Error("LoadWebhook() with corrupted file: expected error, got nil")
	}
}

func TestNewTemplateRef(t *testing.T) {
	ref := NewTemplateRef("check-context-size", "window", nil)

	if ref.Hook != "check-context-size" {
		t.Errorf("Hook = %q, want %q", ref.Hook, "check-context-size")
	}
	if ref.Variant != "window" {
		t.Errorf("Variant = %q, want %q", ref.Variant, "window")
	}
	if ref.Variables != nil {
		t.Errorf("Variables = %v, want nil", ref.Variables)
	}
}

func TestPayload_JSONMarshal(t *testing.T) {
	original := Payload{
		Event:     "loop",
		Message:   "Loop completed",
		SessionID: "sess-42",
		Timestamp: "2026-01-01T00:00:00Z",
		Project:   "myproject",
		Detail: &TemplateRef{
			Hook:      "check-context-size",
			Variant:   "window",
			Variables: map[string]any{"Percentage": 85},
		},
	}

	data, marshalErr := json.Marshal(original)
	if marshalErr != nil {
		t.Fatalf("json.Marshal() error = %v", marshalErr)
	}

	var restored Payload
	if unmarshalErr := json.Unmarshal(data, &restored); unmarshalErr != nil {
		t.Fatalf("json.Unmarshal() error = %v", unmarshalErr)
	}

	if restored.Event != original.Event {
		t.Errorf("Event = %q, want %q", restored.Event, original.Event)
	}
	if restored.Message != original.Message {
		t.Errorf("Message = %q, want %q", restored.Message, original.Message)
	}
	if restored.SessionID != original.SessionID {
		t.Errorf("SessionID = %q, want %q", restored.SessionID, original.SessionID)
	}
	if restored.Timestamp != original.Timestamp {
		t.Errorf("Timestamp = %q, want %q", restored.Timestamp, original.Timestamp)
	}
	if restored.Project != original.Project {
		t.Errorf("Project = %q, want %q", restored.Project, original.Project)
	}
	if restored.Detail == nil {
		t.Fatal("Detail is nil after roundtrip")
	}
	if restored.Detail.Hook != original.Detail.Hook {
		t.Errorf("Detail.Hook = %q, want %q", restored.Detail.Hook, original.Detail.Hook)
	}
	if restored.Detail.Variant != original.Detail.Variant {
		t.Errorf("Detail.Variant = %q, want %q", restored.Detail.Variant, original.Detail.Variant)
	}
}
