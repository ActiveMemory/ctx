//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/spf13/cobra"
)

func newTestCmd() *cobra.Command {
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

func cmdOutput(cmd *cobra.Command) string {
	return cmd.OutOrStdout().(*bytes.Buffer).String()
}

func TestCheckContextSize_SilentEarly(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	// Change to temp dir so .context/logs don't pollute
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-silent"}`)
	if err := runCheckContextSize(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Context Checkpoint") {
		t.Errorf("expected silence at prompt 1, got: %s", out)
	}
}

func TestCheckContextSize_CheckpointAt18(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	// Pre-set counter to 17 so next increment = 18 (18 > 15, 18 is not divisible by 5)
	// Need count 20 for first trigger (20 > 15, 20 % 5 == 0)
	counterFile := filepath.Join(tmpDir, "ctx", "context-check-test-18")
	_ = os.MkdirAll(filepath.Dir(counterFile), 0o700)
	_ = os.WriteFile(counterFile, []byte("19"), 0o600)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-18"}`)
	if err := runCheckContextSize(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Context Checkpoint") {
		t.Errorf("expected checkpoint at prompt 20, got: %s", out)
	}
	if !strings.Contains(out, "prompt #20") {
		t.Errorf("expected 'prompt #20' in output, got: %s", out)
	}
	if !strings.Contains(out, "Context:") {
		t.Errorf("expected context dir footer, got: %s", out)
	}
}

func TestCheckContextSize_CheckpointAt33(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	// Pre-set counter to 32 so next = 33 (33 > 30, 33 % 3 == 0)
	counterFile := filepath.Join(tmpDir, "ctx", "context-check-test-33")
	_ = os.MkdirAll(filepath.Dir(counterFile), 0o700)
	_ = os.WriteFile(counterFile, []byte("32"), 0o600)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-33"}`)
	if err := runCheckContextSize(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Context Checkpoint") {
		t.Errorf("expected checkpoint at prompt 33, got: %s", out)
	}
}

func TestCheckContextSize_OversizeNudgeAtCheckpoint(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	// Create a flag file simulating an oversize injection
	ctxDir := filepath.Join(workDir, config.DirContext)
	stateDir := filepath.Join(ctxDir, config.DirState)
	_ = os.MkdirAll(stateDir, 0o750)
	flagContent := "Context injection oversize warning\n" +
		"===================================\n" +
		"Timestamp: 2026-02-26T14:30:00Z\n" +
		"Injected:  18200 tokens (threshold: 15000)\n\n" +
		"Per-file breakdown:\n" +
		"  CONSTITUTION.md        1200 tokens\n"
	_ = os.WriteFile(filepath.Join(stateDir, "injection-oversize"),
		[]byte(flagContent), 0o600)

	// Set counter to 19 so next = 20 (triggers checkpoint at 20 > 15, 20 % 5 == 0)
	counterFile := filepath.Join(tmpDir, "ctx", "context-check-test-oversize-nudge")
	_ = os.MkdirAll(filepath.Dir(counterFile), 0o700)
	_ = os.WriteFile(counterFile, []byte("19"), 0o600)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-oversize-nudge"}`)
	if err := runCheckContextSize(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Context Checkpoint") {
		t.Error("expected checkpoint output")
	}
	if !strings.Contains(out, "18200") {
		t.Errorf("expected oversize token count in output, got: %s", out)
	}
	if !strings.Contains(out, "ctx-consolidate") {
		t.Error("expected consolidate suggestion in output")
	}

	// Flag should be consumed (deleted)
	flagPath := filepath.Join(stateDir, "injection-oversize")
	if _, err := os.Stat(flagPath); err == nil {
		t.Error("flag file should be deleted after nudge (one-shot)")
	}
}

func TestCheckContextSize_NoFlagNoOversizeNudge(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	// No flag file — trigger a checkpoint
	counterFile := filepath.Join(tmpDir, "ctx", "context-check-test-no-flag")
	_ = os.MkdirAll(filepath.Dir(counterFile), 0o700)
	_ = os.WriteFile(counterFile, []byte("19"), 0o600)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-no-flag"}`)
	if err := runCheckContextSize(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Context Checkpoint") {
		t.Error("expected checkpoint output")
	}
	// Should NOT contain oversize nudge
	if strings.Contains(out, "18200") || strings.Contains(out, "oversize") {
		t.Errorf("should not contain oversize nudge without flag, got: %s", out)
	}
}

func TestCheckContextSize_MalformedFlag(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	// Write a malformed flag file (no parseable token count)
	ctxDir := filepath.Join(workDir, config.DirContext)
	stateDir := filepath.Join(ctxDir, config.DirState)
	_ = os.MkdirAll(stateDir, 0o750)
	_ = os.WriteFile(filepath.Join(stateDir, "injection-oversize"),
		[]byte("garbage data\n"), 0o600)

	counterFile := filepath.Join(tmpDir, "ctx", "context-check-test-malformed")
	_ = os.MkdirAll(filepath.Dir(counterFile), 0o700)
	_ = os.WriteFile(counterFile, []byte("19"), 0o600)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-malformed"}`)
	if err := runCheckContextSize(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	// Should still fire checkpoint, nudge fires with 0 token fallback
	if !strings.Contains(out, "Context Checkpoint") {
		t.Error("expected checkpoint output")
	}
	// Flag should still be consumed
	flagPath := filepath.Join(stateDir, "injection-oversize")
	if _, err := os.Stat(flagPath); err == nil {
		t.Error("malformed flag file should still be consumed")
	}
}

func TestExtractOversizeTokens(t *testing.T) {
	tests := []struct {
		name string
		data string
		want int
	}{
		{
			name: "normal format",
			data: "Injected:  18200 tokens (threshold: 15000)",
			want: 18200,
		},
		{
			name: "single space",
			data: "Injected: 7500 tokens (threshold: 5000)",
			want: 7500,
		},
		{
			name: "no match",
			data: "garbage data",
			want: 0,
		},
		{
			name: "empty",
			data: "",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractOversizeTokens([]byte(tt.data))
			if got != tt.want {
				t.Errorf("extractOversizeTokens() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCheckpointWithTokenLine(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)
	t.Setenv("HOME", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	sessionID := "test-token-line"

	// Create a fake JSONL file with usage data (52k tokens = 26% of 200k)
	projectDir := filepath.Join(tmpDir, ".claude", "projects", "testproj")
	_ = os.MkdirAll(projectDir, 0o750)
	jsonlContent := `{"type":"assistant","message":{"role":"assistant","content":"hi","usage":{"input_tokens":40000,"output_tokens":500,"cache_creation_input_tokens":2000,"cache_read_input_tokens":10000}}}` + "\n"
	_ = os.WriteFile(filepath.Join(projectDir, sessionID+".jsonl"),
		[]byte(jsonlContent), 0o600)

	// Set counter to 19 so next = 20 (triggers checkpoint)
	counterFile := filepath.Join(tmpDir, "ctx", "context-check-"+sessionID)
	_ = os.MkdirAll(filepath.Dir(counterFile), 0o700)
	_ = os.WriteFile(counterFile, []byte("19"), 0o600)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"`+sessionID+`"}`)
	if runErr := runCheckContextSize(cmd, stdin); runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Context Checkpoint") {
		t.Errorf("expected checkpoint, got: %s", out)
	}
	if !strings.Contains(out, "Context window:") {
		t.Errorf("expected token usage line, got: %s", out)
	}
	if !strings.Contains(out, "52k") {
		t.Errorf("expected ~52k tokens in output, got: %s", out)
	}
	// 52k/200k = 26%, should NOT say "running low"
	if strings.Contains(out, "running low") {
		t.Errorf("should not say 'running low' at 26%%, got: %s", out)
	}
}

func TestWindowWarning_Over80(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)
	t.Setenv("HOME", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	sessionID := "test-window-over80"

	// Create a fake JSONL file with 164k tokens (82% of 200k)
	projectDir := filepath.Join(tmpDir, ".claude", "projects", "testproj")
	_ = os.MkdirAll(projectDir, 0o750)
	jsonlContent := `{"type":"assistant","message":{"role":"assistant","content":"hi","usage":{"input_tokens":100000,"output_tokens":2000,"cache_creation_input_tokens":4000,"cache_read_input_tokens":60000}}}` + "\n"
	_ = os.WriteFile(filepath.Join(projectDir, sessionID+".jsonl"),
		[]byte(jsonlContent), 0o600)

	// Counter at 5 — normally silent, but >80% should trigger independently
	counterFile := filepath.Join(tmpDir, "ctx", "context-check-"+sessionID)
	_ = os.MkdirAll(filepath.Dir(counterFile), 0o700)
	_ = os.WriteFile(counterFile, []byte("5"), 0o600)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"`+sessionID+`"}`)
	if runErr := runCheckContextSize(cmd, stdin); runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "Context Window Warning") {
		t.Errorf("expected window warning, got: %s", out)
	}
	if !strings.Contains(out, "82%") {
		t.Errorf("expected 82%% in output, got: %s", out)
	}
}

func TestWindowWarning_Under80_NoCheckpoint(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)
	t.Setenv("HOME", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	sessionID := "test-under80-silent"

	// Create a JSONL file with 40k tokens (20% of 200k)
	projectDir := filepath.Join(tmpDir, ".claude", "projects", "testproj")
	_ = os.MkdirAll(projectDir, 0o750)
	jsonlContent := `{"type":"assistant","message":{"role":"assistant","content":"hi","usage":{"input_tokens":30000,"output_tokens":500,"cache_creation_input_tokens":0,"cache_read_input_tokens":10000}}}` + "\n"
	_ = os.WriteFile(filepath.Join(projectDir, sessionID+".jsonl"),
		[]byte(jsonlContent), 0o600)

	// Counter at 5 — normally silent
	counterFile := filepath.Join(tmpDir, "ctx", "context-check-"+sessionID)
	_ = os.MkdirAll(filepath.Dir(counterFile), 0o700)
	_ = os.WriteFile(counterFile, []byte("5"), 0o600)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"`+sessionID+`"}`)
	if runErr := runCheckContextSize(cmd, stdin); runErr != nil {
		t.Fatalf("unexpected error: %v", runErr)
	}

	out := cmdOutput(cmd)
	if strings.Contains(out, "Context Checkpoint") || strings.Contains(out, "Context Window Warning") {
		t.Errorf("expected silence at prompt 6 with 20%% usage, got: %s", out)
	}
}

func TestTokenUsageLine(t *testing.T) {
	tests := []struct {
		name       string
		tokens     int
		pct        int
		windowSize int
		wantIcon   string
		wantSuffix string
	}{
		{
			name:       "under 80%",
			tokens:     52000,
			pct:        26,
			windowSize: 200000,
			wantIcon:   "⏱",
		},
		{
			name:       "at 80%",
			tokens:     160000,
			pct:        80,
			windowSize: 200000,
			wantIcon:   "⚠",
			wantSuffix: "running low",
		},
		{
			name:       "over 80%",
			tokens:     164000,
			pct:        82,
			windowSize: 200000,
			wantIcon:   "⚠",
			wantSuffix: "running low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tokenUsageLine(tt.tokens, tt.pct, tt.windowSize)
			if !strings.Contains(got, tt.wantIcon) {
				t.Errorf("expected icon %q in %q", tt.wantIcon, got)
			}
			if tt.wantSuffix != "" && !strings.Contains(got, tt.wantSuffix) {
				t.Errorf("expected %q in %q", tt.wantSuffix, got)
			}
			if tt.wantSuffix == "" && strings.Contains(got, "running low") {
				t.Errorf("unexpected 'running low' in %q", got)
			}
		})
	}
}

func TestCheckContextSize_EmptyStdin(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	cmd := newTestCmd()
	stdin := createTempStdin(t, "")
	if err := runCheckContextSize(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should not panic or error with empty input
}

// setupContextDir creates a minimal context directory with essential files so
// that isInitialized() returns true. Must be called after chdir to the work dir.
// Resets rc state so rc.ContextDir() returns the default ".context".
func setupContextDir(t *testing.T) {
	t.Helper()
	rc.Reset()
	dir := rc.ContextDir()
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatal(err)
	}
	for _, f := range config.FilesRequired {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("# "+f+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// createTempStdin writes content to a temp file and returns it opened for reading.
func createTempStdin(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "stdin-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	return f
}
