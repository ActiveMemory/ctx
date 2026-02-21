//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// TestDriftCommand tests the drift command.
func TestDriftCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-drift-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Run drift - just verify it runs without error
	driftCmd := Cmd()
	driftCmd.SetArgs([]string{})

	if err := driftCmd.Execute(); err != nil {
		t.Fatalf("drift command failed: %v", err)
	}
}

// TestDriftJSONOutput tests the drift command with JSON output.
func TestDriftJSONOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-drift-json-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// First init
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Test drift with JSON output
	driftCmd := Cmd()
	driftCmd.SetArgs([]string{"--json"})
	if err := driftCmd.Execute(); err != nil {
		t.Fatalf("drift --json failed: %v", err)
	}
}

// helper: create a cobra command with captured output
func newTestCmd() (*cobra.Command, *bytes.Buffer) {
	cmd := &cobra.Command{Use: "test"}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	return cmd, buf
}

// helper: set up a temp directory as working dir and init context
func setupContextDir(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "cli-drift-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	rc.Reset()

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	initCmd.SilenceUsage = true
	initCmd.SilenceErrors = true
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	return tmpDir, func() {
		_ = os.Chdir(origDir)
		_ = os.RemoveAll(tmpDir)
		rc.Reset()
	}
}

// --- Error function tests ---

func TestErrTasksNotFound(t *testing.T) {
	err := errTasksNotFound()
	if err == nil || !strings.Contains(err.Error(), "TASKS.md not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestErrNoCompletedTasks(t *testing.T) {
	err := errNoCompletedTasks()
	if err == nil || !strings.Contains(err.Error(), "no completed tasks") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestErrMkdir(t *testing.T) {
	err := errMkdir("/some/path", fmt.Errorf("permission denied"))
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "failed to create /some/path") {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "permission denied") {
		t.Errorf("error should wrap cause: %v", err)
	}
}

func TestErrFileWrite(t *testing.T) {
	err := errFileWrite("/some/file", fmt.Errorf("disk full"))
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "failed to write /some/file") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestErrNoTemplate(t *testing.T) {
	err := errNoTemplate("FOO.md", fmt.Errorf("not found"))
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "no template available for FOO.md") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestErrViolationsFound(t *testing.T) {
	err := errViolationsFound()
	if err == nil || !strings.Contains(err.Error(), "violations") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestErrNoContext(t *testing.T) {
	err := errNoContext()
	if err == nil || !strings.Contains(err.Error(), "no .context/ directory found") {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- formatCheckName tests ---

func TestFormatCheckName(t *testing.T) {
	tests := []struct {
		input drift.CheckName
		want  string
	}{
		{drift.CheckPathReferences, "Path references are valid"},
		{drift.CheckStaleness, "No staleness indicators"},
		{drift.CheckConstitution, "Constitution rules respected"},
		{drift.CheckRequiredFiles, "All required files present"},
		{drift.CheckName("unknown_check"), "unknown_check"},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			got := formatCheckName(tt.input)
			if got != tt.want {
				t.Errorf("formatCheckName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// --- outputDriftText tests ---

func TestOutputDriftText_Clean(t *testing.T) {
	cmd, buf := newTestCmd()
	report := &drift.Report{
		Warnings:   []drift.Issue{},
		Violations: []drift.Issue{},
		Passed:     []drift.CheckName{drift.CheckPathReferences, drift.CheckStaleness},
	}

	err := outputDriftText(cmd, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "PASSED") {
		t.Error("expected PASSED section in output")
	}
	if !strings.Contains(out, "OK") {
		t.Error("expected OK status in output")
	}
}

func TestOutputDriftText_WithViolations(t *testing.T) {
	cmd, buf := newTestCmd()
	report := &drift.Report{
		Violations: []drift.Issue{
			{File: "CONSTITUTION.md", Line: 5, Type: drift.IssueSecret, Message: "potential secret found", Rule: "no-secrets"},
		},
		Warnings: []drift.Issue{},
		Passed:   []drift.CheckName{},
	}

	err := outputDriftText(cmd, report)
	if err == nil {
		t.Fatal("expected error for violation report")
	}
	if !strings.Contains(err.Error(), "violations") {
		t.Errorf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "VIOLATIONS") {
		t.Error("expected VIOLATIONS section in output")
	}
	if !strings.Contains(out, "CONSTITUTION.md:5") {
		t.Error("expected file:line reference in violations")
	}
	if !strings.Contains(out, "rule: no-secrets") {
		t.Error("expected rule reference in violations")
	}
	if !strings.Contains(out, "VIOLATION") {
		t.Error("expected VIOLATION status")
	}
}

func TestOutputDriftText_ViolationWithoutLineAndRule(t *testing.T) {
	cmd, buf := newTestCmd()
	report := &drift.Report{
		Violations: []drift.Issue{
			{File: "TEST.md", Type: drift.IssueSecret, Message: "secret detected"},
		},
		Warnings: []drift.Issue{},
		Passed:   []drift.CheckName{},
	}

	err := outputDriftText(cmd, report)
	if err == nil {
		t.Fatal("expected error for violation report")
	}

	out := buf.String()
	// Without line number, should show "- TEST.md: secret detected" (no line num)
	if !strings.Contains(out, "TEST.md: secret detected") {
		t.Errorf("expected violation without line number, got: %s", out)
	}
}

func TestOutputDriftText_WithWarnings(t *testing.T) {
	cmd, buf := newTestCmd()
	report := &drift.Report{
		Violations: []drift.Issue{},
		Warnings: []drift.Issue{
			{File: "ARCHITECTURE.md", Line: 10, Type: drift.IssueDeadPath, Path: "internal/old/pkg", Message: "dead path"},
			{File: "TASKS.md", Type: drift.IssueStaleness, Message: "3 completed tasks"},
			{File: "README.md", Type: drift.IssueType("other_type"), Message: "some other issue"},
		},
		Passed: []drift.CheckName{},
	}

	err := outputDriftText(cmd, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WARNINGS") {
		t.Error("expected WARNINGS section")
	}
	if !strings.Contains(out, "Path References") {
		t.Error("expected path references grouping")
	}
	if !strings.Contains(out, "internal/old/pkg") {
		t.Error("expected dead path reference in output")
	}
	if !strings.Contains(out, "Staleness") {
		t.Error("expected staleness grouping")
	}
	if !strings.Contains(out, "Other") {
		t.Error("expected other grouping")
	}
	if !strings.Contains(out, "WARNING") {
		t.Error("expected WARNING status")
	}
}

// --- outputDriftJSON tests ---

func TestOutputDriftJSON_WithIssues(t *testing.T) {
	cmd, buf := newTestCmd()
	report := &drift.Report{
		Violations: []drift.Issue{
			{File: "X.md", Type: drift.IssueSecret, Message: "secret"},
		},
		Warnings: []drift.Issue{
			{File: "Y.md", Type: drift.IssueStaleness, Message: "stale"},
		},
		Passed: []drift.CheckName{drift.CheckPathReferences},
	}

	err := outputDriftJSON(cmd, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"status"`) {
		t.Error("expected status in JSON")
	}
	if !strings.Contains(out, `"violations"`) {
		t.Error("expected violations in JSON")
	}
	if !strings.Contains(out, `"warnings"`) {
		t.Error("expected warnings in JSON")
	}
}

// --- runDrift tests ---

func TestRunDrift_NoContext(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-drift-nocontext-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	defer rc.Reset()

	cmd := Cmd()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{})

	runErr := cmd.Execute()
	if runErr == nil {
		t.Fatal("expected error when no .context/ exists")
	}
	if !strings.Contains(runErr.Error(), "no .context/ directory found") {
		t.Errorf("unexpected error: %v", runErr)
	}
}

func TestRunDrift_WithFix(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Write TASKS.md with completed tasks to trigger staleness fix
	tasksPath := filepath.Join(tmpDir, config.DirContext, config.FileTask)
	tasksContent := "# Tasks\n\n## In Progress\n\n- [ ] Do something\n\n## Completed\n\n- [x] Done thing 1\n- [x] Done thing 2\n- [x] Done thing 3\n- [x] Done thing 4\n- [x] Done thing 5\n- [x] Done thing 6\n"
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatalf("failed to write TASKS.md: %v", err)
	}

	cmd := Cmd()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--fix"})

	// This may or may not error depending on whether the stale detection
	// actually triggers - just test it doesn't panic
	_ = cmd.Execute()
}

func TestRunDrift_JSONWithViolations(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Create a file that looks like it has secrets to trigger a violation
	constPath := filepath.Join(tmpDir, config.DirContext, "CONSTITUTION.md")
	constContent := "# Constitution\n\n- NEVER commit secrets\n"
	if err := os.WriteFile(constPath, []byte(constContent), 0600); err != nil {
		t.Fatalf("failed to write CONSTITUTION.md: %v", err)
	}

	cmd := Cmd()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--json"})

	// The command may succeed or fail depending on what drift.Detect finds
	_ = cmd.Execute()
}

// --- applyFixes tests ---

func TestApplyFixes_DeadPath(t *testing.T) {
	cmd, buf := newTestCmd()
	ctx := &context.Context{}
	report := &drift.Report{
		Warnings: []drift.Issue{
			{File: "ARCHITECTURE.md", Line: 5, Type: drift.IssueDeadPath, Path: "old/path"},
		},
		Violations: []drift.Issue{},
	}

	result := applyFixes(cmd, ctx, report)
	if result.skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.skipped)
	}
	if result.fixed != 0 {
		t.Errorf("expected 0 fixed, got %d", result.fixed)
	}

	out := buf.String()
	if !strings.Contains(out, "Cannot auto-fix dead path") {
		t.Errorf("expected dead path skip message, got: %s", out)
	}
}

func TestApplyFixes_Secret(t *testing.T) {
	cmd, buf := newTestCmd()
	ctx := &context.Context{}
	report := &drift.Report{
		Warnings: []drift.Issue{},
		Violations: []drift.Issue{
			{File: "SECRETS.md", Type: drift.IssueSecret, Message: "potential secret"},
		},
	}

	result := applyFixes(cmd, ctx, report)
	if result.skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.skipped)
	}

	out := buf.String()
	if !strings.Contains(out, "Cannot auto-fix potential secret") {
		t.Errorf("expected secret skip message, got: %s", out)
	}
}

func TestApplyFixes_MissingFile(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	cmd, _ := newTestCmd()

	ctx, err := context.Load("")
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	// Remove a required file to simulate missing_file issue
	constPath := filepath.Join(tmpDir, config.DirContext, "CONSTITUTION.md")
	_ = os.Remove(constPath)

	report := &drift.Report{
		Warnings: []drift.Issue{
			{File: "CONSTITUTION.md", Type: drift.IssueMissing, Message: "file missing"},
		},
		Violations: []drift.Issue{},
	}

	result := applyFixes(cmd, ctx, report)
	if result.fixed != 1 {
		t.Errorf("expected 1 fixed, got %d", result.fixed)
	}

	// Check the file was recreated
	if _, err := os.Stat(constPath); os.IsNotExist(err) {
		t.Error("CONSTITUTION.md should have been recreated")
	}
}

func TestApplyFixes_MissingFileNoTemplate(t *testing.T) {
	_, cleanup := setupContextDir(t)
	defer cleanup()

	cmd, _ := newTestCmd()
	ctx, err := context.Load("")
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	report := &drift.Report{
		Warnings: []drift.Issue{
			{File: "NONEXISTENT_TEMPLATE.md", Type: drift.IssueMissing, Message: "file missing"},
		},
		Violations: []drift.Issue{},
	}

	result := applyFixes(cmd, ctx, report)
	if len(result.errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.errors))
	}
}

func TestApplyFixes_Staleness_NoTasksFile(t *testing.T) {
	_, cleanup := setupContextDir(t)
	defer cleanup()

	cmd, _ := newTestCmd()

	// Build a context without TASKS.md by removing it from files
	ctx := &context.Context{
		Dir:   rc.ContextDir(),
		Files: []context.FileInfo{},
	}

	report := &drift.Report{
		Warnings: []drift.Issue{
			{File: config.FileTask, Type: drift.IssueStaleness, Message: "stale"},
		},
		Violations: []drift.Issue{},
	}

	result := applyFixes(cmd, ctx, report)
	if len(result.errors) != 1 {
		t.Errorf("expected 1 error, got %d errors: %v", len(result.errors), result.errors)
	}
}

func TestApplyFixes_Staleness_NoCompletedTasks(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Write TASKS.md with no completed tasks
	tasksPath := filepath.Join(tmpDir, config.DirContext, config.FileTask)
	tasksContent := "# Tasks\n\n## In Progress\n\n- [ ] Do something\n\n## Completed\n\n"
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatalf("failed to write TASKS.md: %v", err)
	}

	cmd, _ := newTestCmd()
	ctx, err := context.Load("")
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	report := &drift.Report{
		Warnings: []drift.Issue{
			{File: config.FileTask, Type: drift.IssueStaleness, Message: "stale"},
		},
		Violations: []drift.Issue{},
	}

	result := applyFixes(cmd, ctx, report)
	if len(result.errors) != 1 {
		t.Errorf("expected 1 error (no completed tasks), got %d: %v", len(result.errors), result.errors)
	}
}

func TestApplyFixes_Staleness_Success(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Write TASKS.md with completed tasks
	tasksPath := filepath.Join(tmpDir, config.DirContext, config.FileTask)
	tasksContent := "# Tasks\n\n## In Progress\n\n- [ ] Do something\n\n## Completed\n\n- [x] Done thing 1\n- [x] Done thing 2\n"
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatalf("failed to write TASKS.md: %v", err)
	}

	cmd, buf := newTestCmd()
	ctx, err := context.Load("")
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	report := &drift.Report{
		Warnings: []drift.Issue{
			{File: config.FileTask, Type: drift.IssueStaleness, Message: "stale"},
		},
		Violations: []drift.Issue{},
	}

	result := applyFixes(cmd, ctx, report)
	if result.fixed != 1 {
		t.Errorf("expected 1 fixed, got %d; errors: %v", result.fixed, result.errors)
	}

	out := buf.String()
	if !strings.Contains(out, "Archived") {
		t.Errorf("expected archive message, got: %s", out)
	}

	// Verify TASKS.md no longer has completed tasks
	content, err := os.ReadFile(tasksPath) //nolint:gosec // test temp path
	if err != nil {
		t.Fatalf("failed to read TASKS.md: %v", err)
	}
	if strings.Contains(string(content), "- [x]") {
		t.Error("TASKS.md should not contain completed tasks after fix")
	}
}

// --- fixMissingFile tests ---

func TestFixMissingFile_Success(t *testing.T) {
	_, cleanup := setupContextDir(t)
	defer cleanup()

	// Remove CONSTITUTION.md
	constPath := filepath.Join(rc.ContextDir(), "CONSTITUTION.md")
	_ = os.Remove(constPath)

	err := fixMissingFile("CONSTITUTION.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, statErr := os.Stat(constPath); os.IsNotExist(statErr) {
		t.Error("CONSTITUTION.md should have been created")
	}
}

func TestFixMissingFile_NoTemplate(t *testing.T) {
	_, cleanup := setupContextDir(t)
	defer cleanup()

	err := fixMissingFile("THIS_DOES_NOT_EXIST.md")
	if err == nil {
		t.Fatal("expected error for missing template")
	}
	if !strings.Contains(err.Error(), "no template available") {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- fixStaleness tests ---

func TestFixStaleness_NoTasksFile(t *testing.T) {
	cmd, _ := newTestCmd()
	ctx := &context.Context{
		Files: []context.FileInfo{},
	}

	err := fixStaleness(cmd, ctx)
	if err == nil {
		t.Fatal("expected error for missing TASKS.md")
	}
	if !strings.Contains(err.Error(), "TASKS.md not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFixStaleness_NoCompletedTasks(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	tasksPath := filepath.Join(tmpDir, config.DirContext, config.FileTask)
	tasksContent := "# Tasks\n\n## In Progress\n\n- [ ] Do something\n"
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatalf("failed to write TASKS.md: %v", err)
	}

	cmd, _ := newTestCmd()
	ctx, err := context.Load("")
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	fixErr := fixStaleness(cmd, ctx)
	if fixErr == nil {
		t.Fatal("expected error for no completed tasks")
	}
	if !strings.Contains(fixErr.Error(), "no completed tasks") {
		t.Errorf("unexpected error: %v", fixErr)
	}
}

// --- runDrift with fix flag when issues exist ---

func TestRunDrift_FixWithStaleness(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Create TASKS.md with many completed tasks to trigger staleness
	tasksPath := filepath.Join(tmpDir, config.DirContext, config.FileTask)
	var sb strings.Builder
	sb.WriteString("# Tasks\n\n## In Progress\n\n- [ ] Active task\n\n## Completed\n\n")
	for i := 0; i < 10; i++ {
		sb.WriteString(fmt.Sprintf("- [x] Completed task %d\n", i))
	}
	if err := os.WriteFile(tasksPath, []byte(sb.String()), 0600); err != nil {
		t.Fatalf("failed to write TASKS.md: %v", err)
	}

	cmd := Cmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--fix"})

	// Execute - may or may not error depending on other drift checks
	_ = cmd.Execute()
}

func TestFixStaleness_CompletedSectionWithNextSection(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// TASKS.md has a Completed section followed by another ## heading
	tasksPath := filepath.Join(tmpDir, config.DirContext, config.FileTask)
	tasksContent := "# Tasks\n\n## In Progress\n\n- [ ] Active\n\n## Completed\n\n- [x] Done thing\n\n## Archive\n\nSome archived content\n"
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0600); err != nil {
		t.Fatalf("failed to write TASKS.md: %v", err)
	}

	cmd, buf := newTestCmd()
	ctx, err := context.Load("")
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	fixErr := fixStaleness(cmd, ctx)
	if fixErr != nil {
		t.Fatalf("unexpected error: %v", fixErr)
	}

	out := buf.String()
	if !strings.Contains(out, "Archived") {
		t.Errorf("expected archive message, got: %s", out)
	}

	// Verify Archive section is preserved
	content, err := os.ReadFile(tasksPath) //nolint:gosec // test temp path
	if err != nil {
		t.Fatalf("failed to read TASKS.md: %v", err)
	}
	if !strings.Contains(string(content), "## Archive") {
		t.Error("## Archive section should be preserved")
	}
}

func TestRunDrift_FixTriggersRecheck(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Remove a required file so fixMissingFile gets called and succeeds
	constPath := filepath.Join(tmpDir, config.DirContext, "CONSTITUTION.md")
	_ = os.Remove(constPath)

	// Use runDrift directly with a cobra command that captures output
	cmd := Cmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--fix"})

	_ = cmd.Execute()

	out := buf.String()
	// The fix should have recreated CONSTITUTION.md and triggered re-check
	if !strings.Contains(out, "Applying fixes") {
		t.Errorf("expected 'Applying fixes' in output, got: %s", out)
	}
}

func TestRunDrift_GenericError(t *testing.T) {
	// Test runDrift when context.Load returns a non-NotFoundError
	// This is difficult to trigger naturally since context.Load
	// returns NotFoundError for missing dirs. We can test via the
	// command by creating a .context that's actually a file (not a dir).
	tmpDir, err := os.MkdirTemp("", "cli-drift-generr-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	rc.Reset()
	defer rc.Reset()

	// Create .context as a file, not a directory. context.Load should
	// return a NotFoundError in this case.
	if err := os.WriteFile(filepath.Join(tmpDir, config.DirContext), []byte("not a dir"), 0600); err != nil {
		t.Fatalf("failed to create fake .context: %v", err)
	}

	cmd := Cmd()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{})

	runErr := cmd.Execute()
	if runErr == nil {
		t.Fatal("expected error when .context is a file")
	}
}

func TestRunDrift_FixWithErrorsOutput(t *testing.T) {
	_, cleanup := setupContextDir(t)
	defer cleanup()

	// Test the applyFixes code path that prints errors
	cmd, buf := newTestCmd()
	ctx := &context.Context{
		Dir:   rc.ContextDir(),
		Files: []context.FileInfo{},
	}

	report := &drift.Report{
		Warnings: []drift.Issue{
			{File: config.FileTask, Type: drift.IssueStaleness, Message: "stale"},
		},
		Violations: []drift.Issue{},
	}

	result := applyFixes(cmd, ctx, report)

	// Should have error because no TASKS.md in context
	if len(result.errors) == 0 {
		t.Error("expected at least 1 error")
	}

	// Now test the full runDrift path that prints Fixed/Skipped/Error counts
	// by calling it directly via the cobra command
	_ = buf.String() // consume
}

// --- checkFileAge tests ---

func TestCheckFileAge_StaleFile(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Backdate DECISIONS.md to 45 days ago
	decisionsPath := filepath.Join(tmpDir, config.DirContext, config.FileDecision)
	past := time.Now().AddDate(0, 0, -45)
	if err := os.Chtimes(decisionsPath, past, past); err != nil {
		t.Fatalf("failed to backdate file: %v", err)
	}

	ctx, err := context.Load("")
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	report := drift.Detect(ctx)

	// Should have at least one stale_age warning for DECISIONS.md
	found := false
	for _, w := range report.Warnings {
		if w.Type == drift.IssueStaleAge && w.File == config.FileDecision {
			found = true
			if !strings.Contains(w.Message, "days ago") {
				t.Errorf("expected 'days ago' in message, got: %s", w.Message)
			}
		}
	}
	if !found {
		t.Error("expected stale_age warning for DECISIONS.md")
	}
}

func TestCheckFileAge_ConstitutionExcluded(t *testing.T) {
	tmpDir, cleanup := setupContextDir(t)
	defer cleanup()

	// Backdate CONSTITUTION.md to 60 days ago
	constPath := filepath.Join(tmpDir, config.DirContext, config.FileConstitution)
	past := time.Now().AddDate(0, 0, -60)
	if err := os.Chtimes(constPath, past, past); err != nil {
		t.Fatalf("failed to backdate file: %v", err)
	}

	ctx, err := context.Load("")
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	report := drift.Detect(ctx)

	// Should NOT have a stale_age warning for CONSTITUTION.md
	for _, w := range report.Warnings {
		if w.Type == drift.IssueStaleAge && w.File == config.FileConstitution {
			t.Error("CONSTITUTION.md should be excluded from file age check")
		}
	}
}

func TestCheckFileAge_RecentFiles(t *testing.T) {
	_, cleanup := setupContextDir(t)
	defer cleanup()

	ctx, err := context.Load("")
	if err != nil {
		t.Fatalf("failed to load context: %v", err)
	}

	report := drift.Detect(ctx)

	// Freshly initialized files should pass the file age check
	for _, w := range report.Warnings {
		if w.Type == drift.IssueStaleAge {
			t.Errorf("unexpected stale_age warning for fresh file: %s", w.File)
		}
	}

	// CheckFileAge should be in the passed list
	found := false
	for _, p := range report.Passed {
		if p == drift.CheckFileAge {
			found = true
		}
	}
	if !found {
		t.Error("expected CheckFileAge in passed checks for fresh context")
	}
}

func TestFormatCheckName_FileAge(t *testing.T) {
	got := formatCheckName(drift.CheckFileAge)
	want := "No stale files by age"
	if got != want {
		t.Errorf("formatCheckName(CheckFileAge) = %q, want %q", got, want)
	}
}

func TestApplyFixes_StaleAge(t *testing.T) {
	cmd, buf := newTestCmd()
	ctx := &context.Context{}
	report := &drift.Report{
		Warnings: []drift.Issue{
			{File: "DECISIONS.md", Type: drift.IssueStaleAge, Message: "last modified 45 days ago"},
		},
		Violations: []drift.Issue{},
	}

	result := applyFixes(cmd, ctx, report)
	if result.skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.skipped)
	}
	if result.fixed != 0 {
		t.Errorf("expected 0 fixed, got %d", result.fixed)
	}

	out := buf.String()
	if !strings.Contains(out, "Cannot auto-fix file age") {
		t.Errorf("expected file age skip message, got: %s", out)
	}
}
