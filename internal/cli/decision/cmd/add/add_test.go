//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"os"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

// TestDecisionAdd verifies the noun-first ctx decision add
// subcommand writes a structured ADR entry to DECISIONS.md.
func TestDecisionAdd(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-decision-add-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	testctx.Declare(t, tmpDir)

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	addCmd := Cmd()
	addCmd.SetArgs([]string{
		"Use PostgreSQL for database",
		"--session-id", "test1234",
		"--branch", "main",
		"--commit", "abc123",
		"--context", "Need a reliable database",
		"--rationale", "PostgreSQL is well-supported",
		"--consequence", "Team needs training",
	})
	if err = addCmd.Execute(); err != nil {
		t.Fatalf("ctx decision add failed: %v", err)
	}

	content, err := os.ReadFile(".context/DECISIONS.md")
	if err != nil {
		t.Fatalf("failed to read DECISIONS.md: %v", err)
	}
	contentStr := string(content)
	for _, want := range []string{
		"Use PostgreSQL for database",
		"Need a reliable database",
		"PostgreSQL is well-supported",
		"Team needs training",
	} {
		if !strings.Contains(contentStr, want) {
			t.Errorf("expected %q in DECISIONS.md", want)
		}
	}
}

// TestDecisionAddRequiresFlags verifies that omitting
// required flags produces an error.
func TestDecisionAddRequiresFlags(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-decision-add-req-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	testctx.Declare(t, tmpDir)

	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	addCmd := Cmd()
	addCmd.SetArgs([]string{"Incomplete decision"})
	err = addCmd.Execute()
	if err == nil {
		t.Fatal("expected error when adding decision without required flags")
	}
	if !strings.Contains(err.Error(), "--session-id") {
		t.Errorf("error should mention missing --session-id flag: %v", err)
	}
}
