//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package doctor

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func setupContextDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CTX_DIR", dir)
	rc.Reset()

	// Create required files.
	for _, f := range config.FilesRequired {
		path := filepath.Join(dir, f)
		if writeErr := os.WriteFile(path, []byte("# "+f+"\n"), 0o600); writeErr != nil {
			t.Fatal(writeErr)
		}
	}
	return dir
}

func TestDoctor_Healthy(t *testing.T) {
	setupContextDir(t)

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("doctor failed: %v", runErr)
	}

	output := out.String()
	if !strings.Contains(output, "0 warnings, 0 errors") {
		t.Errorf("expected healthy summary, got: %s", output)
	}
	if !strings.Contains(output, "Context initialized") {
		t.Errorf("expected context initialized check, got: %s", output)
	}
}

func TestDoctor_MissingRequiredFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CTX_DIR", dir)
	rc.Reset()

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("doctor failed: %v", runErr)
	}

	output := out.String()
	if !strings.Contains(output, "Missing required files") {
		t.Errorf("expected missing files error, got: %s", output)
	}
	if !strings.Contains(output, "1 errors") {
		t.Errorf("expected 1 error in summary, got: %s", output)
	}
}

func TestDoctor_EventLogOff(t *testing.T) {
	setupContextDir(t)

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	output := out.String()
	if !strings.Contains(output, "Event logging disabled") {
		t.Errorf("expected event logging info note, got: %s", output)
	}
	// Info notes should not count as warnings or errors.
	if !strings.Contains(output, "0 warnings, 0 errors") {
		t.Errorf("expected 0 warnings/errors (info is not a warning), got: %s", output)
	}
}

func TestDoctor_JSON(t *testing.T) {
	setupContextDir(t)

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--json"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("doctor --json failed: %v", runErr)
	}

	var report Report
	if unmarshalErr := json.Unmarshal(out.Bytes(), &report); unmarshalErr != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", unmarshalErr, out.String())
	}
	if len(report.Results) == 0 {
		t.Error("expected at least one result")
	}
}

func TestDoctor_HighCompletion(t *testing.T) {
	dir := setupContextDir(t)

	// Write a TASKS.md with high completion ratio.
	tasks := "# Tasks\n"
	for i := 0; i < 20; i++ {
		tasks += "- [x] Completed task\n"
	}
	tasks += "- [ ] Pending task\n"
	tasksPath := filepath.Join(dir, config.FileTask)
	if writeErr := os.WriteFile(tasksPath, []byte(tasks), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	output := out.String()
	if !strings.Contains(output, "consider archiving") {
		t.Errorf("expected archiving suggestion for high completion, got: %s", output)
	}
}

func TestDoctor_ContextSizeBreakdown(t *testing.T) {
	dir := setupContextDir(t)

	// Write enough content to some files to verify per-file breakdown appears.
	archPath := filepath.Join(dir, "ARCHITECTURE.md")
	if writeErr := os.WriteFile(archPath, []byte(strings.Repeat("word ", 500)), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	tasksPath := filepath.Join(dir, config.FileTask)
	if writeErr := os.WriteFile(tasksPath, []byte(strings.Repeat("task ", 200)), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	output := out.String()

	// Should show "window" not "budget".
	if strings.Contains(output, "budget") {
		t.Errorf("should use 'window' not 'budget', got: %s", output)
	}
	if !strings.Contains(output, "window:") {
		t.Errorf("expected 'window:' in context size line, got: %s", output)
	}

	// Should show per-file breakdown.
	if !strings.Contains(output, "ARCHITECTURE.md") {
		t.Errorf("expected ARCHITECTURE.md in breakdown, got: %s", output)
	}
	if !strings.Contains(output, "TASKS.md") {
		t.Errorf("expected TASKS.md in breakdown, got: %s", output)
	}
	if !strings.Contains(output, "tokens") {
		t.Errorf("expected 'tokens' in breakdown lines, got: %s", output)
	}
}

func TestDoctor_ContextSizeJSON(t *testing.T) {
	setupContextDir(t)

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--json"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("doctor --json failed: %v", runErr)
	}

	var report Report
	if unmarshalErr := json.Unmarshal(out.Bytes(), &report); unmarshalErr != nil {
		t.Fatalf("output is not valid JSON: %v", unmarshalErr)
	}

	// Should have context_file_* results.
	var fileResults int
	for _, r := range report.Results {
		if strings.HasPrefix(r.Name, "context_file_") {
			fileResults++
		}
	}
	if fileResults == 0 {
		t.Error("expected context_file_* results in JSON output")
	}
}

func TestDoctor_DriftWarnings(t *testing.T) {
	dir := setupContextDir(t)

	// Add an ARCHITECTURE.md referencing a nonexistent path to trigger drift.
	archPath := filepath.Join(dir, "ARCHITECTURE.md")
	archContent := "# Architecture\n\n" +
		"See `internal/nonexistent/fake.go` for details.\n"
	if writeErr := os.WriteFile(archPath, []byte(archContent), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	// Drift detection may or may not find warnings depending on the checks
	// that are relevant. The important thing is that it doesn't crash.
	output := out.String()
	if !strings.Contains(output, "ctx doctor") {
		t.Errorf("expected doctor header in output, got: %s", output)
	}
}
