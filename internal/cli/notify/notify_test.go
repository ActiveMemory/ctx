//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func setupCLITest(t *testing.T) (string, func()) {
	t.Helper()
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	_ = os.MkdirAll(filepath.Join(tempDir, ".context"), 0o750)
	// Create required files so isInitialized returns true
	for _, f := range config.FilesRequired {
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

	err = runSetup(cmd, tmpFile)
	if err != nil {
		t.Fatalf("runSetup() error = %v", err)
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
		got := maskURL(tc.input)
		if got != tc.want {
			t.Errorf("maskURL(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
