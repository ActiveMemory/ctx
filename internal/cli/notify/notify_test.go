//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/notify/cmd/setup"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	notifylib "github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func setupCLITest(t *testing.T) (string, func()) {
	t.Helper()
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	_ = os.MkdirAll(filepath.Join(tempDir, ".context"), 0o750)
	// Create required files so isInitialized returns true
	for _, f := range ctx.FilesRequired {
		_ = os.WriteFile(filepath.Join(tempDir, ".context", f), []byte("# "+f+"\n"), 0o600)
	}
	rc.Reset()
	return tempDir, func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	}
}

func TestCmd_MissingEvent(t *testing.T) {
	_, cleanup := setupCLITest(t)
	defer cleanup()

	cmd := Cmd()
	cmd.SetArgs([]string{"hello"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --event")
	}
	if !strings.Contains(err.Error(), "event") {
		t.Errorf("error = %q, want mention of 'event'", err.Error())
	}
}

func TestCmd_MissingMessage(t *testing.T) {
	_, cleanup := setupCLITest(t)
	defer cleanup()

	cmd := Cmd()
	cmd.SetArgs([]string{"--event", "test"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing message")
	}
	if !strings.Contains(err.Error(), "message") {
		t.Errorf("error = %q, want mention of 'message'", err.Error())
	}
}

func TestCmd_NoopNoWebhook(t *testing.T) {
	_, cleanup := setupCLITest(t)
	defer cleanup()

	cmd := Cmd()
	cmd.SetArgs([]string{"--event", "test", "hello from test"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestSetup_WithMockStdin(t *testing.T) {
	_, cleanup := setupCLITest(t)
	defer cleanup()

	// Create a temp file to use as stdin
	tmpFile, err := os.CreateTemp("", "notify-stdin-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	_, _ = tmpFile.WriteString("https://example.com/webhook?key=secret\n")
	_, _ = tmpFile.Seek(0, 0)

	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err = setup.Run(cmd, tmpFile)
	if err != nil {
		t.Fatalf("setup.Run() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Webhook configured") {
		t.Errorf("output = %q, want 'Webhook configured'", output)
	}
	if strings.Contains(output, "secret") {
		t.Error("output should not contain the full URL secret")
	}
}

func TestMaskURL(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"https://example.com/webhook?key=secret", "https://example.com/***"},
		{"https://hooks.slack.com/services/T00/B00/xxx", "https://hooks.slack.com/***"},
		{"http://localhost:8080", "http://localhost:808***"},
	}

	for _, tc := range tests {
		got := notifylib.MaskURL(tc.input)
		if got != tc.want {
			t.Errorf("notifylib.MaskURL(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestSetup_EmptyInput(t *testing.T) {
	_, cleanup := setupCLITest(t)
	defer cleanup()

	tmpFile, createErr := os.CreateTemp("", "notify-stdin-empty-*")
	if createErr != nil {
		t.Fatal(createErr)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	// Write only a newline so the scanner reads an empty line.
	_, _ = tmpFile.WriteString("\n")
	_, _ = tmpFile.Seek(0, 0)

	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	setupErr := setup.Run(cmd, tmpFile)
	if setupErr == nil {
		t.Fatal("expected error for empty webhook URL input")
	}
	if !strings.Contains(setupErr.Error(), "empty") {
		t.Errorf("error = %q, want mention of 'empty'", setupErr.Error())
	}
}

func TestTest_NoWebhookConfigured(t *testing.T) {
	_, cleanup := setupCLITest(t)
	defer cleanup()

	cmd := Cmd()
	cmd.SetArgs([]string{"test"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	execErr := cmd.Execute()
	if execErr != nil {
		t.Fatalf("Execute() error = %v", execErr)
	}

	output := buf.String()
	if !strings.Contains(strings.ToLower(output), "no webhook") {
		t.Errorf("output = %q, want mention of 'no webhook'", output)
	}
}

func TestTest_WebhookSuccess(t *testing.T) {
	_, cleanup := setupCLITest(t)
	defer cleanup()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Save the test server URL as the webhook.
	if saveErr := notifylib.SaveWebhook(server.URL); saveErr != nil {
		t.Fatalf("SaveWebhook() error = %v", saveErr)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"test"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	execErr := cmd.Execute()
	if execErr != nil {
		t.Fatalf("Execute() error = %v", execErr)
	}

	output := buf.String()
	if !strings.Contains(output, "200") {
		t.Errorf("output = %q, want mention of HTTP 200", output)
	}
	if !strings.Contains(output, "working") {
		t.Errorf("output = %q, want mention of 'working'", output)
	}
}

func TestTest_WebhookServerError(t *testing.T) {
	_, cleanup := setupCLITest(t)
	defer cleanup()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	if saveErr := notifylib.SaveWebhook(server.URL); saveErr != nil {
		t.Fatalf("SaveWebhook() error = %v", saveErr)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"test"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	execErr := cmd.Execute()
	if execErr != nil {
		t.Fatalf("Execute() error = %v", execErr)
	}

	output := buf.String()
	if !strings.Contains(output, "500") {
		t.Errorf("output = %q, want mention of HTTP 500", output)
	}
}
