//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package remind

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// setup creates a temp dir with a .context/ directory and sets the RC override.
func setup(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	rc.Reset()
	rc.OverrideContextDir(config.DirContext)

	ctxDir := filepath.Join(dir, config.DirContext)
	if err := os.MkdirAll(ctxDir, 0750); err != nil {
		t.Fatal(err)
	}

	return dir
}

// runCmd executes a cobra command and captures its output.
func runCmd(cmd *cobra.Command) (string, error) {
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	return buf.String(), err
}

// newRemindCmd builds a fresh remind command with the given args.
func newRemindCmd(args ...string) *cobra.Command {
	cmd := Cmd()
	cmd.SetArgs(args)
	return cmd
}

func TestAdd_Basic(t *testing.T) {
	setup(t)

	out, err := runCmd(newRemindCmd("add", "refactor the swagger definitions"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "+ [1] refactor the swagger definitions") {
		t.Errorf("output = %q, want confirmation with ID 1", out)
	}

	// Verify JSON file content.
	data, err := os.ReadFile(RemindersPath())
	if err != nil {
		t.Fatalf("read reminders file: %v", err)
	}
	var reminders []Reminder
	if err := json.Unmarshal(data, &reminders); err != nil {
		t.Fatalf("parse reminders: %v", err)
	}
	if len(reminders) != 1 {
		t.Fatalf("got %d reminders, want 1", len(reminders))
	}
	if reminders[0].ID != 1 {
		t.Errorf("ID = %d, want 1", reminders[0].ID)
	}
	if reminders[0].Message != "refactor the swagger definitions" {
		t.Errorf("Message = %q, want %q", reminders[0].Message, "refactor the swagger definitions")
	}
	if reminders[0].After != nil {
		t.Errorf("After = %v, want nil", reminders[0].After)
	}
}

func TestAdd_WithAfter(t *testing.T) {
	setup(t)

	out, err := runCmd(newRemindCmd("add", "check CI", "--after", "2099-01-15"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "(after 2099-01-15)") {
		t.Errorf("output = %q, want date annotation", out)
	}

	reminders, err := ReadReminders()
	if err != nil {
		t.Fatalf("read reminders: %v", err)
	}
	if reminders[0].After == nil {
		t.Fatal("After is nil, want date string")
	}
	if *reminders[0].After != "2099-01-15" {
		t.Errorf("After = %q, want %q", *reminders[0].After, "2099-01-15")
	}
}

func TestAdd_InvalidDate(t *testing.T) {
	setup(t)

	_, err := runCmd(newRemindCmd("add", "test", "--after", "garbage"))
	if err == nil {
		t.Fatal("expected error for invalid date")
	}
	if !strings.Contains(err.Error(), "invalid date") {
		t.Errorf("error = %q, want 'invalid date'", err.Error())
	}
}

func TestAdd_IDIncrement(t *testing.T) {
	setup(t)

	// Add three reminders.
	_, _ = runCmd(newRemindCmd("add", "first"))
	_, _ = runCmd(newRemindCmd("add", "second"))
	_, _ = runCmd(newRemindCmd("add", "third"))

	// Dismiss the middle one (ID 2).
	_, _ = runCmd(newRemindCmd("dismiss", "2"))

	// Add another — should get ID 4 (max existing is 3, so 3+1).
	out, err := runCmd(newRemindCmd("add", "fourth"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "+ [4]") {
		t.Errorf("output = %q, want ID 4", out)
	}
}

func TestAdd_DefaultAction(t *testing.T) {
	setup(t)

	// "ctx remind TEXT" without "add" subcommand should work.
	out, err := runCmd(newRemindCmd("refactor swagger"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "+ [1] refactor swagger") {
		t.Errorf("output = %q, want confirmation", out)
	}
}

func TestList_Empty(t *testing.T) {
	setup(t)

	out, err := runCmd(newRemindCmd("list"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No reminders.") {
		t.Errorf("output = %q, want 'No reminders.'", out)
	}
}

func TestList_Mixed(t *testing.T) {
	setup(t)

	// Add a due reminder (no date gate).
	_, _ = runCmd(newRemindCmd("add", "due now"))
	// Add a not-yet-due reminder.
	_, _ = runCmd(newRemindCmd("add", "future task", "--after", "2099-12-31"))

	out, err := runCmd(newRemindCmd("list"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "[1] due now") {
		t.Errorf("output missing due reminder: %q", out)
	}
	if !strings.Contains(out, "not yet due") {
		t.Errorf("output missing 'not yet due' annotation: %q", out)
	}
}

func TestDismiss_ByID(t *testing.T) {
	setup(t)

	_, _ = runCmd(newRemindCmd("add", "to dismiss"))
	_, _ = runCmd(newRemindCmd("add", "to keep"))

	out, err := runCmd(newRemindCmd("dismiss", "1"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "- [1] to dismiss") {
		t.Errorf("output = %q, want dismissal confirmation", out)
	}

	// Verify only one remains.
	reminders, err := ReadReminders()
	if err != nil {
		t.Fatalf("read reminders: %v", err)
	}
	if len(reminders) != 1 {
		t.Fatalf("got %d reminders, want 1", len(reminders))
	}
	if reminders[0].Message != "to keep" {
		t.Errorf("remaining = %q, want 'to keep'", reminders[0].Message)
	}
}

func TestDismiss_NotFound(t *testing.T) {
	setup(t)

	_, err := runCmd(newRemindCmd("dismiss", "99"))
	if err == nil {
		t.Fatal("expected error for nonexistent ID")
	}
	if !strings.Contains(err.Error(), "no reminder with ID 99") {
		t.Errorf("error = %q, want to contain 'no reminder with ID 99'", err.Error())
	}
}

func TestDismiss_All(t *testing.T) {
	setup(t)

	_, _ = runCmd(newRemindCmd("add", "first"))
	_, _ = runCmd(newRemindCmd("add", "second"))

	out, err := runCmd(newRemindCmd("dismiss", "--all"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Dismissed 2 reminders.") {
		t.Errorf("output = %q, want 'Dismissed 2 reminders.'", out)
	}

	// Verify file is empty array.
	reminders, err := ReadReminders()
	if err != nil {
		t.Fatalf("read reminders: %v", err)
	}
	if len(reminders) != 0 {
		t.Errorf("got %d reminders, want 0", len(reminders))
	}
}
