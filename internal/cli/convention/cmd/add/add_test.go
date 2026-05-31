//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

// TestConventionAdd verifies the noun-first ctx convention add
// subcommand writes an entry to CONVENTIONS.md.
func TestConventionAdd(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-convention-add-*")
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
	addCmd.SetArgs([]string{"Use camelCase for variable names"})
	if err = addCmd.Execute(); err != nil {
		t.Fatalf("ctx convention add failed: %v", err)
	}

	content, err := os.ReadFile(".context/CONVENTIONS.md")
	if err != nil {
		t.Fatalf("failed to read CONVENTIONS.md: %v", err)
	}
	if !strings.Contains(string(content), "Use camelCase for variable names") {
		t.Error("convention was not added to CONVENTIONS.md")
	}
}

// TestConventionAddFromJSONFile verifies --json-file supplies the
// convention's content via the title field (convention has no other
// structured fields).
func TestConventionAddFromJSONFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-convention-add-json-*")
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

	payload := filepath.Join(tmpDir, "convention.json")
	if err = os.WriteFile(payload, []byte(
		`{"title": "Resolve binaries from /usr/local/bin on PATH"}`,
	), 0o600); err != nil {
		t.Fatalf("write payload: %v", err)
	}

	addCmd := Cmd()
	addCmd.SetArgs([]string{"--json-file", payload})
	if err = addCmd.Execute(); err != nil {
		t.Fatalf("ctx convention add --json-file failed: %v", err)
	}

	content, err := os.ReadFile(".context/CONVENTIONS.md")
	if err != nil {
		t.Fatalf("failed to read CONVENTIONS.md: %v", err)
	}
	if !strings.Contains(
		string(content), "Resolve binaries from /usr/local/bin on PATH",
	) {
		t.Error("convention content from --json-file was not added")
	}
}
