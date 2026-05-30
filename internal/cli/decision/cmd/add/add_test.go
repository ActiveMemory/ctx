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

// TestDecisionAddFromJSONFile verifies that --json-file populates the
// typed fields and provenance from a JSON envelope, including a
// rationale value whose content would trip a literal command-string
// deny rule (the feature's driver).
func TestDecisionAddFromJSONFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-decision-add-json-*")
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

	payload := filepath.Join(tmpDir, "decision.json")
	if err = os.WriteFile(payload, []byte(`{
		"title": "Install ctx into the system PATH",
		"context": "agents invoke ctx by bare name",
		"rationale": "the binary belongs at /usr/local/bin so it is on PATH",
		"consequence": "ctx resolves from any working directory",
		"provenance": {"session_id": "json1234", "branch": "main", "commit": "abc123"}
	}`), 0o600); err != nil {
		t.Fatalf("write payload: %v", err)
	}

	addCmd := Cmd()
	addCmd.SetArgs([]string{"--json-file", payload})
	if err = addCmd.Execute(); err != nil {
		t.Fatalf("ctx decision add --json-file failed: %v", err)
	}

	content, err := os.ReadFile(".context/DECISIONS.md")
	if err != nil {
		t.Fatalf("failed to read DECISIONS.md: %v", err)
	}
	contentStr := string(content)
	for _, want := range []string{
		"Install ctx into the system PATH",
		"agents invoke ctx by bare name",
		"the binary belongs at /usr/local/bin so it is on PATH",
		"ctx resolves from any working directory",
	} {
		if !strings.Contains(contentStr, want) {
			t.Errorf("expected %q in DECISIONS.md", want)
		}
	}
}

// TestDecisionAddJSONFileRejectsPlaceholder verifies the placeholder
// gate also fires on JSON-supplied values, so --json-file is not a
// bypass for the schema checks.
func TestDecisionAddJSONFileRejectsPlaceholder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-decision-add-json-ph-*")
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

	payload := filepath.Join(tmpDir, "decision.json")
	if err = os.WriteFile(payload, []byte(`{
		"title": "Decision body",
		"context": "real context",
		"rationale": "TBD",
		"consequence": "real consequence",
		"provenance": {"session_id": "json1234", "branch": "main", "commit": "abc123"}
	}`), 0o600); err != nil {
		t.Fatalf("write payload: %v", err)
	}

	addCmd := Cmd()
	addCmd.SetArgs([]string{"--json-file", payload})
	err = addCmd.Execute()
	if err == nil {
		t.Fatal("expected placeholder rejection for JSON rationale=TBD")
	}
	if !strings.Contains(err.Error(), "rationale") {
		t.Errorf("error should name --rationale: %v", err)
	}
}

// TestDecisionAddRequiresFlags verifies that omitting
// required body flags produces an error. The placeholder
// check fires first (PreRunE runs before cobra's required-
// flag validation), so the message names the first empty
// body flag rather than --session-id.
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
	if !strings.Contains(err.Error(), "--context") {
		t.Errorf("error should mention missing --context flag: %v", err)
	}
}

// TestDecisionAddRejectsPlaceholderRationale verifies that
// passing a placeholder body-flag value (TBD, see chat, etc.)
// fails fast at PreRunE.
func TestDecisionAddRejectsPlaceholderRationale(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-decision-add-ph-*")
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
		"Decision body",
		"--session-id", "test1234",
		"--branch", "main",
		"--commit", "abc123",
		"--context", "real context",
		"--rationale", "TBD",
		"--consequence", "real consequence",
	})
	err = addCmd.Execute()
	if err == nil {
		t.Fatal("expected placeholder rejection for --rationale=TBD")
	}
	if !strings.Contains(err.Error(), "rationale") {
		t.Errorf("error should name --rationale: %v", err)
	}
	if !strings.Contains(err.Error(), "placeholder") {
		t.Errorf("error should explain placeholder rejection: %v", err)
	}
}

// TestDecisionAddRejectsWhitespaceContext verifies whitespace-
// only values are rejected at PreRunE.
func TestDecisionAddRejectsWhitespaceContext(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-decision-add-ws-*")
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
		"Decision body",
		"--session-id", "test1234",
		"--branch", "main",
		"--commit", "abc123",
		"--context", "   ",
		"--rationale", "real rationale",
		"--consequence", "real consequence",
	})
	err = addCmd.Execute()
	if err == nil {
		t.Fatal("expected rejection for whitespace --context")
	}
	if !strings.Contains(err.Error(), "context") {
		t.Errorf("error should name --context: %v", err)
	}
}
