//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/crypto"
	"github.com/ActiveMemory/ctx/internal/entity"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	_ = os.MkdirAll(filepath.Join(tempDir, ".context"), 0o750)

	testctx.Declare(t, tempDir)

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

	// Create the (global) key but no encrypted file. The key resolves
	// to ~/.ctx/.ctx.key — HOME is redirected to tempDir in tests.
	globalDir := filepath.Join(tempDir, ".ctx")
	if err := os.MkdirAll(globalDir, 0o700); err != nil {
		t.Fatal(err)
	}
	keyPath := filepath.Join(globalDir, crypto.ContextKey)
	if err := os.WriteFile(keyPath, make([]byte, 32), 0o600); err != nil {
		t.Fatal(err)
	}

	url, err := LoadWebhook()
	if err != nil {
		t.Fatalf("LoadWebhook() error = %v", err)
	}
	if url != "" {
		t.Errorf("LoadWebhook() = %q, want empty", url)
	}
}

func TestLoadWebhook_InvalidKeyPropagated(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Key present (wrong size) AND the encrypted file present: the
	// webhook IS configured, so an invalid key is a real
	// misconfiguration LoadWebhook must surface, not a silent "no
	// webhook". (Before the swallow was narrowed it returned ("", nil)
	// and the failure vanished.)
	globalDir := filepath.Join(tempDir, ".ctx")
	if err := os.MkdirAll(globalDir, 0o700); err != nil {
		t.Fatal(err)
	}
	keyPath := filepath.Join(globalDir, crypto.ContextKey)
	if err := os.WriteFile(keyPath, []byte("too-short"), 0o600); err != nil {
		t.Fatal(err)
	}
	encPath := filepath.Join(tempDir, ".context", crypto.NotifyEnc)
	if err := os.WriteFile(encPath, []byte("ciphertext"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadWebhook(); err == nil {
		t.Fatal("LoadWebhook() with an invalid-size key: expected error, got nil")
	}
}

func TestLoadWebhook_ConfiguredKeyAbsentSurfaces(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	// .notify.enc present (webhook IS configured) but no key anywhere.
	// This is "configured but broken", not "not configured", so
	// LoadWebhook must surface an error rather than silently report no
	// webhook (the absent-key-in-worktree / fresh-machine case).
	encPath := filepath.Join(tempDir, ".context", crypto.NotifyEnc)
	if err := os.WriteFile(encPath, []byte("ciphertext"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadWebhook(); err == nil {
		t.Fatal("LoadWebhook() with enc present but key absent: expected error, got nil")
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
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Subscribe the event so Send passes the event filter and actually
	// reaches the webhook-absence check, rather than short-circuiting
	// at the filter. No .notify.enc exists, so this is the legitimate
	// "not configured" path: noop with no error and no warning.
	rcContent := "notify:\n  events:\n    - test\n"
	if err := os.WriteFile(
		filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0o600,
	); err != nil {
		t.Fatal(err)
	}
	rc.Reset()

	var buf bytes.Buffer
	defer logWarn.SetSink(&buf)()

	if err := Send("test", "hello", "session-1", nil); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("unconfigured webhook should not warn, got %q", buf.String())
	}
}

func TestSend_EventFiltered(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Set up a server that should NOT be called
	called := false
	ts := httptest.NewServer(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			called = true
		}),
	)
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
	ts := httptest.NewServer(
		http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&received)
		}),
	)
	defer ts.Close()

	if err := SaveWebhook(ts.URL); err != nil {
		t.Fatalf("SaveWebhook() error = %v", err)
	}

	// Configure events to allow "loop"
	rcContent := "notify:\n  events:\n    - loop\n"
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0o600)
	rc.Reset()

	ref := entity.NewTemplateRef("check-context-size", "window",
		map[string]any{"Percentage": 82, "TokenCount": "164k"})
	sendErr := Send("loop", "Loop completed after 5 iterations", "abc123", ref)
	if sendErr != nil {
		t.Fatalf("Send() error = %v", sendErr)
	}

	if received["event"] != "loop" {
		t.Errorf("Event = %v, want %q", received["event"], "loop")
	}
	if received["message"] != "Loop completed after 5 iterations" {
		t.Errorf(
			"Message = %v, want %q",
			received["message"],
			"Loop completed after 5 iterations",
		)
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
		t.Fatalf(
			"Detail is not an object: %T = %v",
			received["detail"], received["detail"],
		)
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
		t.Errorf(
			"Detail.variables.TokenCount = %v, want %q",
			vars["TokenCount"], "164k",
		)
	}
}

func TestSend_NilDetailOmitted(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	var received map[string]any
	ts := httptest.NewServer(
		http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&received)
		}),
	)
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

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer ts.Close()

	if err := SaveWebhook(ts.URL); err != nil {
		t.Fatalf("SaveWebhook() error = %v", err)
	}

	// Configure events to allow "test"
	rcContent := "notify:\n  events:\n    - test\n"
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0o600)
	rc.Reset()

	// An HTTP 500 is a received response, not a transport error, so
	// this exercises the success/Body.Close path. The transport-error
	// (postErr) warn branch is covered by TestSend_PostFailureWarns.
	err := Send("test", "hello", "session-1", nil)
	if err != nil {
		t.Fatalf("Send() error = %v, want nil (fire-and-forget)", err)
	}
}

func TestSend_ConfiguredButUndeliverableWarns(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	called := false
	ts := httptest.NewServer(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			called = true
		}),
	)
	defer ts.Close()

	// Configure a webhook, then make it undeliverable by replacing the
	// key with a different valid key so .notify.enc no longer decrypts
	// — the worktree / wrong-key footgun. Send must warn (not silently
	// drop) and must not POST.
	if err := SaveWebhook(ts.URL); err != nil {
		t.Fatalf("SaveWebhook() error = %v", err)
	}
	keyPath := filepath.Join(tempDir, ".ctx", crypto.ContextKey)
	if err := os.WriteFile(
		keyPath, bytes.Repeat([]byte{7}, 32), 0o600,
	); err != nil {
		t.Fatal(err)
	}
	rcContent := "notify:\n  events:\n    - stop\n"
	if err := os.WriteFile(
		filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0o600,
	); err != nil {
		t.Fatal(err)
	}
	rc.Reset()

	var buf bytes.Buffer
	defer logWarn.SetSink(&buf)()

	if err := Send("stop", "hello", "s1", nil); err != nil {
		t.Fatalf("Send() error = %v, want nil (fire-and-forget)", err)
	}
	if called {
		t.Error("webhook was POSTed despite a decrypt failure")
	}
	if n := strings.Count(
		buf.String(), "webhook configured but undeliverable",
	); n != 1 {
		t.Errorf("undeliverable warning count = %d, want 1; sink=%q",
			n, buf.String())
	}
}

func TestSend_PostFailureWarns(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Stand up a server, capture its URL, then close it so the POST
	// gets connection-refused (a transport error, unlike an HTTP 500).
	ts := httptest.NewServer(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	url := ts.URL
	ts.Close()

	if err := SaveWebhook(url); err != nil {
		t.Fatalf("SaveWebhook() error = %v", err)
	}
	rcContent := "notify:\n  events:\n    - stop\n"
	if err := os.WriteFile(
		filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0o600,
	); err != nil {
		t.Fatal(err)
	}
	rc.Reset()

	var buf bytes.Buffer
	defer logWarn.SetSink(&buf)()

	if err := Send("stop", "hello", "s1", nil); err != nil {
		t.Fatalf("Send() error = %v, want nil (fire-and-forget)", err)
	}
	if n := strings.Count(buf.String(), "webhook POST failed"); n != 1 {
		t.Errorf("POST-failure warning count = %d, want 1; sink=%q",
			n, buf.String())
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
	if writeErr := os.WriteFile(
		encPath, []byte("corrupted-garbage-data"), 0o600,
	); writeErr != nil {
		t.Fatalf("WriteFile() error = %v", writeErr)
	}

	// Should return an error, not panic.
	_, loadErr := LoadWebhook()
	if loadErr == nil {
		t.Error("LoadWebhook() with corrupted file: expected error, got nil")
	}
}

func TestNewTemplateRef(t *testing.T) {
	ref := entity.NewTemplateRef("check-context-size", "window", nil)

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
	original := entity.NotifyPayload{
		Event:     "loop",
		Message:   "Loop completed",
		SessionID: "sess-42",
		Timestamp: "2026-01-01T00:00:00Z",
		Project:   "myproject",
		Detail: &entity.TemplateRef{
			Hook:      "check-context-size",
			Variant:   "window",
			Variables: map[string]any{"Percentage": 85},
		},
	}

	data, marshalErr := json.Marshal(original)
	if marshalErr != nil {
		t.Fatalf("json.Marshal() error = %v", marshalErr)
	}

	var restored entity.NotifyPayload
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
		t.Errorf(
			"Detail.Hook = %q, want %q", restored.Detail.Hook, original.Detail.Hook,
		)
	}
	if restored.Detail.Variant != original.Detail.Variant {
		t.Errorf(
			"Detail.Variant = %q, want %q",
			restored.Detail.Variant, original.Detail.Variant,
		)
	}
}
