//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
	"github.com/ActiveMemory/ctx/internal/trace"
)

func TestTraceTagAndShow(t *testing.T) {
	tmpDir := t.TempDir()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	testctx.Declare(t, tmpDir)

	// Init git repo
	run(t, "git", "init")
	run(t, "git", "config", "user.email", "test@test.com")
	run(t, "git", "config", "user.name", "Test")

	// Init ctx
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init: %v", err)
	}

	// Create a file and commit
	if err := os.WriteFile("test.go", []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run(t, "git", "add", ".")
	run(t, "git", "commit", "-m", "Initial commit")

	// Record some pending context
	stateDir := filepath.Join(".context", "state")
	if err := os.MkdirAll(stateDir, 0750); err != nil {
		t.Fatal(err)
	}
	_ = trace.Record("decision:1", stateDir)

	// Write history for the commit
	traceDir := filepath.Join(".context", "trace")
	hash := strings.TrimSpace(runOutput(t, "git", "rev-parse", "HEAD"))

	histErr := trace.WriteHistory(trace.HistoryEntry{
		Commit:  hash,
		Refs:    []string{"decision:1"},
		Message: "Initial commit",
	}, traceDir)
	if histErr != nil {
		t.Fatalf("WriteHistory: %v", histErr)
	}

	// Test ctx trace <commit>; should not error
	showCmd := Cmd()
	showCmd.SetArgs([]string{hash[:7]})
	if showErr := showCmd.Execute(); showErr != nil {
		t.Errorf("trace show failed: %v", showErr)
	}

	// Test ctx trace --last 5
	lastCmd := Cmd()
	lastCmd.SetArgs([]string{"--last", "5"})
	if lastErr := lastCmd.Execute(); lastErr != nil {
		t.Errorf("trace --last failed: %v", lastErr)
	}

	// Test ctx trace tag
	tagCmd := Cmd()
	tagCmd.SetArgs([]string{"tag", "HEAD", "--note", "Test tag"})
	if tagErr := tagCmd.Execute(); tagErr != nil {
		t.Errorf("trace tag failed: %v", tagErr)
	}

	// Verify override was written
	overrides, ovrErr := trace.ReadOverrides(traceDir)
	if ovrErr != nil {
		t.Fatalf("ReadOverrides: %v", ovrErr)
	}
	if len(overrides) != 1 {
		t.Errorf("expected 1 override, got %d", len(overrides))
	}
}

func run(t *testing.T, name string, args ...string) {
	t.Helper()
	//nolint:gosec // test helper, name is always "git" from test code
	cmd := exec.Command(name, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
	}
}

func runOutput(t *testing.T, name string, args ...string) string {
	t.Helper()
	//nolint:gosec // test helper, name is always "git" from test code
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		t.Fatalf("%s %v failed: %v", name, args, err)
	}
	return string(out)
}
