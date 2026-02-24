//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// makeEntries creates n entries with valid ## [YYYY-MM-DD-HHMMSS] headers.
func makeEntries(n int) string {
	var b strings.Builder
	b.WriteString("# Heading\n\n")
	for i := range n {
		b.WriteString("## [2026-01-15-")
		b.WriteString(strings.Repeat("0", 5))
		b.WriteString(string(rune('0' + i%10)))
		b.WriteString("] Entry ")
		b.WriteString(string(rune('A' + i%26)))
		b.WriteString("\n\nBody.\n\n")
	}
	return b.String()
}

// makeLines creates a conventions-style file with n lines.
func makeLines(n int) string {
	var b strings.Builder
	b.WriteString("# Conventions\n\n")
	for i := range n {
		b.WriteString(fmt.Sprintf("- Convention %d: do something\n", i))
	}
	return b.String()
}

func setupKnowledgeTest(t *testing.T, decEntries, lrnEntries int) string {
	t.Helper()
	return setupKnowledgeTestFull(t, decEntries, lrnEntries, -1)
}

func setupKnowledgeTestFull(t *testing.T, decEntries, lrnEntries, convLines int) string {
	t.Helper()
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	_ = os.Chdir(dir)

	rc.Reset()
	t.Cleanup(rc.Reset)

	ctxDir := filepath.Join(dir, config.DirContext)
	_ = os.MkdirAll(ctxDir, 0o750)

	// Create required files so isInitialized() returns true
	for _, f := range config.FilesRequired {
		_ = os.WriteFile(filepath.Join(ctxDir, f), []byte("# "+f+"\n"), 0o600)
	}

	// Overwrite with entry content
	if decEntries >= 0 {
		_ = os.WriteFile(filepath.Join(ctxDir, config.FileDecision),
			[]byte(makeEntries(decEntries)), 0o600)
	}
	if lrnEntries >= 0 {
		_ = os.WriteFile(filepath.Join(ctxDir, config.FileLearning),
			[]byte(makeEntries(lrnEntries)), 0o600)
	}
	if convLines >= 0 {
		_ = os.WriteFile(filepath.Join(ctxDir, config.FileConvention),
			[]byte(makeLines(convLines)), 0o600)
	}

	// Point XDG_RUNTIME_DIR to a temp location so marker files don't interfere
	tmpState := filepath.Join(dir, "state")
	_ = os.MkdirAll(tmpState, 0o700)
	t.Setenv("XDG_RUNTIME_DIR", tmpState)

	return dir
}

func TestCheckKnowledge_BelowThreshold(t *testing.T) {
	setupKnowledgeTest(t, 5, 5)

	cmd := checkKnowledgeCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected silent output, got:\n%s", buf.String())
	}
}

func TestCheckKnowledge_AboveThreshold(t *testing.T) {
	// Default thresholds: decisions=20, learnings=30
	setupKnowledgeTest(t, 25, 35)

	cmd := checkKnowledgeCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "VERBATIM") {
		t.Error("expected VERBATIM relay header")
	}
	if !strings.Contains(out, "DECISIONS.md") {
		t.Error("expected DECISIONS.md mentioned")
	}
	if !strings.Contains(out, "LEARNINGS.md") {
		t.Error("expected LEARNINGS.md mentioned")
	}
	if !strings.Contains(out, "25") {
		t.Error("expected decision count 25 in output")
	}
	if !strings.Contains(out, "35") {
		t.Error("expected learning count 35 in output")
	}
}

func TestCheckKnowledge_DailyThrottle(t *testing.T) {
	setupKnowledgeTest(t, 25, 35)

	cmd1 := checkKnowledgeCmd()
	var buf1 bytes.Buffer
	cmd1.SetOut(&buf1)
	_ = cmd1.RunE(cmd1, nil)

	if buf1.Len() == 0 {
		t.Fatal("first call should produce output")
	}

	// Second call same day should be silent
	cmd2 := checkKnowledgeCmd()
	var buf2 bytes.Buffer
	cmd2.SetOut(&buf2)
	_ = cmd2.RunE(cmd2, nil)

	if buf2.Len() != 0 {
		t.Errorf("second call should be silent (throttled), got:\n%s", buf2.String())
	}
}

func TestCheckKnowledge_OneOverOneUnder(t *testing.T) {
	// Decisions over (25 > 20), learnings under (10 < 30)
	setupKnowledgeTest(t, 25, 10)

	cmd := checkKnowledgeCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DECISIONS.md") {
		t.Error("expected DECISIONS.md mentioned")
	}
	if strings.Contains(out, "LEARNINGS.md") {
		t.Error("LEARNINGS.md should NOT be mentioned (under threshold)")
	}
}

func TestCheckKnowledge_ConventionsOverLineCount(t *testing.T) {
	// Default threshold: 200 lines. Create 250 lines, entries below threshold.
	setupKnowledgeTestFull(t, 5, 5, 250)

	cmd := checkKnowledgeCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "CONVENTIONS.md") {
		t.Error("expected CONVENTIONS.md mentioned")
	}
	if !strings.Contains(out, "lines") {
		t.Error("expected 'lines' unit for conventions")
	}
	if strings.Contains(out, "DECISIONS.md") {
		t.Error("DECISIONS.md should NOT be mentioned (under threshold)")
	}
}

func TestCheckKnowledge_ConventionsBelowThreshold(t *testing.T) {
	// 100 lines, well below default 200
	setupKnowledgeTestFull(t, 5, 5, 100)

	cmd := checkKnowledgeCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected silent output, got:\n%s", buf.String())
	}
}

func TestCheckKnowledge_Uninitialized(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	_ = os.Chdir(dir)

	rc.Reset()
	t.Cleanup(rc.Reset)

	cmd := checkKnowledgeCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected silent when uninitialized, got:\n%s", buf.String())
	}
}
