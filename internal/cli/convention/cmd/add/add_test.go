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
	// --section is required for conventions (M4/T16): a convention is a
	// bullet under an H2 section, and the CLI refuses to pick one.
	addCmd.SetArgs([]string{
		"--section", "Naming", "Use camelCase for variable names",
	})
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
	// The JSON payload carries the content; the section is still a flag,
	// and still required (M4/T16).
	addCmd.SetArgs([]string{
		"--section", "CLI Structure", "--json-file", payload,
	})
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

// TestConventionAddRequiresSection pins the M4/T16 contract: a
// convention must name the H2 section it belongs under. There is no
// default — a catch-all section is where an undecided caller dumps
// everything, and it defeats the grouping the digesting pass folds on.
// Placeholder values are refused for the same reason.
func TestConventionAddRequiresSection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-convention-add-section-*")
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

	before, err := os.ReadFile(".context/CONVENTIONS.md")
	if err != nil {
		t.Fatalf("read CONVENTIONS.md: %v", err)
	}

	cases := []struct {
		name string
		args []string
	}{
		{"omitted", []string{"Prefer errors.Is at boundaries"}},
		{"empty", []string{"--section", "", "Prefer errors.Is"}},
		{"whitespace only", []string{"--section", "   ", "Prefer errors.Is"}},
		{"placeholder TBD", []string{"--section", "TBD", "Prefer errors.Is"}},
		{"placeholder lowercase", []string{"--section", "tbd", "Prefer errors.Is"}},
		{"placeholder n/a", []string{"--section", "n/a", "Prefer errors.Is"}},
		{"placeholder pending", []string{"--section", "pending", "Prefer errors.Is"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := Cmd()
			cmd.SetArgs(tc.args)
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			if execErr := cmd.Execute(); execErr == nil {
				t.Fatal("Execute err = nil, want a section refusal")
			}
			after, readErr := os.ReadFile(".context/CONVENTIONS.md")
			if readErr != nil {
				t.Fatalf("read CONVENTIONS.md: %v", readErr)
			}
			if string(after) != string(before) {
				t.Error("CONVENTIONS.md mutated by a refused add")
			}
		})
	}
}
