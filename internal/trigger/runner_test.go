//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	cfgTrigger "github.com/ActiveMemory/ctx/internal/config/trigger"
)

// writeHookScript creates an executable shell script in the given hook
// type subdirectory. It creates the directory structure if needed.
func writeHookScript(t *testing.T, hooksDir, hookType, name, body string) {
	t.Helper()
	typeDir := filepath.Join(hooksDir, hookType)
	if err := os.MkdirAll(typeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(typeDir, name)
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

// skipIfWindows skips the test on Windows where shell scripts cannot run.
func skipIfWindows(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("shell script execution not supported on Windows")
	}
}

// TestRunAll_CancelPropagation verifies that when a hook returns
// cancel:true, subsequent hooks are not executed and the aggregated
// output reflects the cancellation.
// Validates: Requirements 7.3
func TestRunAll_CancelPropagation(t *testing.T) {
	skipIfWindows(t)

	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	// First hook (alphabetically) cancels.
	writeHookScript(t, hooksDir, "pre-tool-use", "01-block.sh",
		"#!/bin/sh\necho '{\"cancel\": true, \"message\": \"blocked\"}'")

	// Second hook should never run; if it did it would add context.
	writeHookScript(t, hooksDir, "pre-tool-use", "02-context.sh",
		"#!/bin/sh\necho '{\"cancel\": false, \"context\": \"should not appear\"}'")

	input := &HookInput{TriggerType: "pre-tool-use", Tool: "test"}
	agg, err := RunAll(hooksDir, cfgTrigger.PreToolUse, input, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !agg.Cancelled {
		t.Error("expected Cancelled to be true")
	}
	if agg.Message != "blocked" {
		t.Errorf("expected message %q, got %q", "blocked", agg.Message)
	}
	if agg.Context != "" {
		t.Errorf("expected empty context (second hook should not run), got %q", agg.Context)
	}
}

// TestRunAll_ContextAggregation verifies that non-empty context fields
// from multiple hooks are concatenated with newlines.
// Validates: Requirements 7.4
func TestRunAll_ContextAggregation(t *testing.T) {
	skipIfWindows(t)

	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	writeHookScript(t, hooksDir, "session-start", "01-first.sh",
		"#!/bin/sh\necho '{\"cancel\": false, \"context\": \"extra context\"}'")

	writeHookScript(t, hooksDir, "session-start", "02-second.sh",
		"#!/bin/sh\necho '{\"cancel\": false, \"context\": \"more context\"}'")

	input := &HookInput{TriggerType: "session-start", Tool: "test"}
	agg, err := RunAll(hooksDir, cfgTrigger.SessionStart, input, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if agg.Cancelled {
		t.Error("expected Cancelled to be false")
	}

	want := "extra context\nmore context"
	if agg.Context != want {
		t.Errorf("expected context %q, got %q", want, agg.Context)
	}
}

// TestRunAll_NonZeroExitCode verifies that a hook exiting with a
// non-zero exit code is logged, skipped, and remaining hooks continue.
// Validates: Requirements 7.5
func TestRunAll_NonZeroExitCode(t *testing.T) {
	skipIfWindows(t)

	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	// First hook exits with error.
	writeHookScript(t, hooksDir, "post-tool-use", "01-fail.sh",
		"#!/bin/sh\nexit 1")

	// Second hook succeeds with context.
	writeHookScript(t, hooksDir, "post-tool-use", "02-ok.sh",
		"#!/bin/sh\necho '{\"cancel\": false, \"context\": \"survived\"}'")

	input := &HookInput{TriggerType: "post-tool-use", Tool: "test"}
	agg, err := RunAll(hooksDir, cfgTrigger.PostToolUse, input, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if agg.Cancelled {
		t.Error("expected Cancelled to be false")
	}
	if len(agg.Errors) == 0 {
		t.Error("expected at least one error from the failing hook")
	}
	if agg.Context != "survived" {
		t.Errorf("expected context %q, got %q", "survived", agg.Context)
	}
}

// TestRunAll_InvalidJSONOutput verifies that a hook producing invalid
// JSON on stdout is logged, skipped, and remaining hooks continue.
// Validates: Requirements 7.6
func TestRunAll_InvalidJSONOutput(t *testing.T) {
	skipIfWindows(t)

	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	// First hook outputs invalid JSON.
	writeHookScript(t, hooksDir, "file-save", "01-bad.sh",
		"#!/bin/sh\necho 'not json'")

	// Second hook succeeds.
	writeHookScript(t, hooksDir, "file-save", "02-good.sh",
		"#!/bin/sh\necho '{\"cancel\": false, \"context\": \"valid\"}'")

	input := &HookInput{TriggerType: "file-save", Tool: "test"}
	agg, err := RunAll(hooksDir, cfgTrigger.FileSave, input, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if agg.Cancelled {
		t.Error("expected Cancelled to be false")
	}
	if len(agg.Errors) == 0 {
		t.Error("expected at least one error from the invalid JSON hook")
	}

	// Verify the error mentions invalid JSON.
	found := false
	for _, e := range agg.Errors {
		if strings.Contains(e, "invalid JSON") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected an error containing %q, got %v", "invalid JSON", agg.Errors)
	}

	if agg.Context != "valid" {
		t.Errorf("expected context %q, got %q", "valid", agg.Context)
	}
}

// TestRunAll_TimeoutEnforcement verifies that a hook exceeding the
// timeout is terminated and remaining hooks continue.
// Validates: Requirements 7.7, 7.8, 19.6
func TestRunAll_TimeoutEnforcement(t *testing.T) {
	skipIfWindows(t)

	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	// First hook sleeps well beyond the timeout.
	// Use exec to replace the shell process with sleep so that
	// CommandContext's kill signal reaches the sleeping process
	// directly, avoiding orphaned child process issues.
	writeHookScript(t, hooksDir, "context-add", "01-slow.sh",
		"#!/bin/sh\nexec sleep 30")

	// Second hook succeeds quickly.
	writeHookScript(t, hooksDir, "context-add", "02-fast.sh",
		"#!/bin/sh\necho '{\"cancel\": false, \"context\": \"fast\"}'")

	input := &HookInput{TriggerType: "context-add", Tool: "test"}

	start := time.Now()
	agg, err := RunAll(hooksDir, cfgTrigger.ContextAdd, input, 1*time.Second)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should complete well under 30s (the sleep duration).
	if elapsed > 10*time.Second {
		t.Errorf("expected timeout enforcement to complete quickly, took %v", elapsed)
	}

	if agg.Cancelled {
		t.Error("expected Cancelled to be false")
	}
	if len(agg.Errors) == 0 {
		t.Error("expected at least one error from the timed-out hook")
	}

	// Verify the error mentions timeout.
	found := false
	for _, e := range agg.Errors {
		if strings.Contains(e, "timeout") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected an error containing %q, got %v", "timeout", agg.Errors)
	}

	if agg.Context != "fast" {
		t.Errorf("expected context %q from second hook, got %q", "fast", agg.Context)
	}
}

// TestRunAll_NoHooks verifies that RunAll returns an empty
// AggregatedOutput when no hooks exist for the given type.
// Validates: Requirements 7.1
func TestRunAll_NoHooks(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	// Create the hooks directory but no type subdirectories.
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatal(err)
	}

	input := &HookInput{TriggerType: "session-end", Tool: "test"}
	agg, err := RunAll(hooksDir, cfgTrigger.SessionEnd, input, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if agg.Cancelled {
		t.Error("expected Cancelled to be false")
	}
	if agg.Context != "" {
		t.Errorf("expected empty context, got %q", agg.Context)
	}
	if len(agg.Errors) != 0 {
		t.Errorf("expected no errors, got %v", agg.Errors)
	}
}

// TestRunAll_MissingHooksDir verifies that RunAll returns an empty
// AggregatedOutput when the hooks directory does not exist.
// Validates: Requirements 7.1
func TestRunAll_MissingHooksDir(t *testing.T) {
	input := &HookInput{TriggerType: "pre-tool-use", Tool: "test"}
	agg, err := RunAll(filepath.Join(t.TempDir(), "nonexistent"), cfgTrigger.PreToolUse, input, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if agg.Cancelled {
		t.Error("expected Cancelled to be false")
	}
	if agg.Context != "" {
		t.Errorf("expected empty context, got %q", agg.Context)
	}
}

// TestRunAll_EmptyStdout verifies that a hook producing no output
// is treated as a no-op (no cancel, no context, no error).
// Validates: Requirements 7.2
func TestRunAll_EmptyStdout(t *testing.T) {
	skipIfWindows(t)

	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	writeHookScript(t, hooksDir, "session-end", "01-silent.sh",
		"#!/bin/sh\n# produces no output")

	input := &HookInput{TriggerType: "session-end", Tool: "test"}
	agg, err := RunAll(hooksDir, cfgTrigger.SessionEnd, input, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if agg.Cancelled {
		t.Error("expected Cancelled to be false")
	}
	if agg.Context != "" {
		t.Errorf("expected empty context, got %q", agg.Context)
	}
	if len(agg.Errors) != 0 {
		t.Errorf("expected no errors, got %v", agg.Errors)
	}
}
