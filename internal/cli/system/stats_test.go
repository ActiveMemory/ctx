//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestStats_NoFiles(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	cmd := newTestCmd()
	cmd.SetArgs([]string{})
	if err := runStats(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := cmdOutput(cmd)
	if !strings.Contains(out, "No stats recorded") {
		t.Errorf("expected 'No stats recorded', got: %s", out)
	}
}

func TestStats_DumpEntries(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	sd := filepath.Join(rc.ContextDir(), config.DirState)
	content := `{"ts":"2026-03-05T10:00:00Z","prompt":1,"tokens":5000,"pct":2,"window":200000,"model":"claude-opus-4-6","event":"silent"}
{"ts":"2026-03-05T10:01:00Z","prompt":2,"tokens":12000,"pct":6,"window":200000,"model":"claude-opus-4-6","event":"silent"}
{"ts":"2026-03-05T10:02:00Z","prompt":3,"tokens":20000,"pct":10,"window":200000,"model":"claude-opus-4-6","event":"checkpoint"}
`
	_ = os.WriteFile(filepath.Join(sd, "stats-abc12345.jsonl"), []byte(content), 0o600)

	cmd := newTestCmd()
	cmd.Flags().BoolP("follow", "f", false, "")
	cmd.Flags().StringP("session", "s", "", "")
	cmd.Flags().IntP("last", "n", 20, "")
	cmd.Flags().BoolP("json", "j", false, "")

	if err := runStats(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "abc12345") {
		t.Errorf("expected session ID in output, got: %s", out)
	}
	if !strings.Contains(out, "checkpoint") {
		t.Errorf("expected 'checkpoint' event in output, got: %s", out)
	}
	if !strings.Contains(out, "PROMPT") {
		t.Errorf("expected header row, got: %s", out)
	}
}

func TestStats_JSONOutput(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	sd := filepath.Join(rc.ContextDir(), config.DirState)
	content := `{"ts":"2026-03-05T10:00:00Z","prompt":1,"tokens":5000,"pct":2,"window":200000,"event":"silent"}
`
	_ = os.WriteFile(filepath.Join(sd, "stats-sess1.jsonl"), []byte(content), 0o600)

	cmd := newTestCmd()
	cmd.Flags().BoolP("follow", "f", false, "")
	cmd.Flags().StringP("session", "s", "", "")
	cmd.Flags().IntP("last", "n", 20, "")
	cmd.Flags().BoolP("json", "j", false, "")
	_ = cmd.Flags().Set("json", "true")

	if err := runStats(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, `"session":"sess1"`) {
		t.Errorf("expected session field in JSON output, got: %s", out)
	}
}

func TestStats_SessionFilter(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	sd := filepath.Join(rc.ContextDir(), config.DirState)
	_ = os.WriteFile(filepath.Join(sd, "stats-match123.jsonl"),
		[]byte(`{"ts":"2026-03-05T10:00:00Z","prompt":1,"tokens":5000,"pct":2,"window":200000,"event":"silent"}`+"\n"), 0o600)
	_ = os.WriteFile(filepath.Join(sd, "stats-other456.jsonl"),
		[]byte(`{"ts":"2026-03-05T10:00:00Z","prompt":1,"tokens":9000,"pct":4,"window":200000,"event":"silent"}`+"\n"), 0o600)

	cmd := newTestCmd()
	cmd.Flags().BoolP("follow", "f", false, "")
	cmd.Flags().StringP("session", "s", "", "")
	cmd.Flags().IntP("last", "n", 20, "")
	cmd.Flags().BoolP("json", "j", false, "")
	_ = cmd.Flags().Set("session", "match")

	if err := runStats(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "match123") {
		t.Errorf("expected matched session, got: %s", out)
	}
	if strings.Contains(out, "other456") {
		t.Errorf("filtered session should not appear, got: %s", out)
	}
}

func TestStats_LastN(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	sd := filepath.Join(rc.ContextDir(), config.DirState)
	var lines string
	for i := 1; i <= 10; i++ {
		lines += `{"ts":"2026-03-05T10:00:0` + string(rune('0'+i%10)) + `Z","prompt":` +
			strings.TrimSpace(strings.Repeat(" ", 0)) +
			`1,"tokens":5000,"pct":2,"window":200000,"event":"silent"}` + "\n"
	}
	_ = os.WriteFile(filepath.Join(sd, "stats-lastn.jsonl"), []byte(lines), 0o600)

	cmd := newTestCmd()
	cmd.Flags().BoolP("follow", "f", false, "")
	cmd.Flags().StringP("session", "s", "", "")
	cmd.Flags().IntP("last", "n", 20, "")
	cmd.Flags().BoolP("json", "j", false, "")
	_ = cmd.Flags().Set("last", "3")
	_ = cmd.Flags().Set("json", "true")

	if err := runStats(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := strings.TrimSpace(cmdOutput(cmd))
	lineCount := len(strings.Split(out, "\n"))
	if lineCount != 3 {
		t.Errorf("expected 3 lines with --last 3, got %d: %s", lineCount, out)
	}
}

func TestExtractSessionID(t *testing.T) {
	tests := []struct {
		basename string
		want     string
	}{
		{"stats-abc123.jsonl", "abc123"},
		{"stats-long-uuid-value.jsonl", "long-uuid-value"},
		{"stats-.jsonl", ""},
	}
	for _, tt := range tests {
		got := extractSessionID(tt.basename)
		if got != tt.want {
			t.Errorf("extractSessionID(%q) = %q, want %q", tt.basename, got, tt.want)
		}
	}
}

func TestReadNewLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stats-test.jsonl")
	first := `{"ts":"2026-03-05T10:00:00Z","prompt":1,"tokens":5000,"pct":2,"window":200000,"event":"silent"}` + "\n"
	second := `{"ts":"2026-03-05T10:01:00Z","prompt":2,"tokens":12000,"pct":6,"window":200000,"event":"checkpoint"}` + "\n"

	_ = os.WriteFile(path, []byte(first), 0o600)
	offset := int64(len(first))

	// Append second line.
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o600)
	_, _ = f.WriteString(second)
	_ = f.Close()

	entries := readNewLines(path, offset, "test")
	if len(entries) != 1 {
		t.Fatalf("expected 1 new entry, got %d", len(entries))
	}
	if entries[0].Event != "checkpoint" {
		t.Errorf("expected 'checkpoint' event, got %q", entries[0].Event)
	}
}
