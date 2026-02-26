//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

func TestContextLoadGate_NotInitialized(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	// No .context/ — should be silent
	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-no-init"}`)
	if err := runContextLoadGate(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("expected silence when not initialized, got: %s", out)
	}
}

func TestContextLoadGate_EmptySessionID(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{}`)
	if err := runContextLoadGate(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if out != "" {
		t.Errorf("expected silence with empty session_id, got: %s", out)
	}
}

func TestContextLoadGate_InjectsContent(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	// Write test content to context files
	ctxDir := filepath.Join(workDir, config.DirContext)
	writeTestFile(t, ctxDir, config.FileConstitution,
		"# Constitution\n\nNever commit secrets.\n")
	writeTestFile(t, ctxDir, config.FileConvention,
		"# Conventions\n\nUse filepath.Join.\n")
	writeTestFile(t, ctxDir, config.FileArchitecture,
		"# Architecture\n\nCLI tool with .context/ directory.\n")
	writeTestFile(t, ctxDir, config.FileAgentPlaybook,
		"# Agent Playbook\n\nWork → Reflect → Persist.\n")
	writeTestFile(t, ctxDir, config.FileDecision,
		"# Decisions\n\n"+
			config.IndexStart+"\n"+
			"| Date | Decision |\n"+
			"|------|----------|\n"+
			"| 2026-02-26 | Use auto-injection |\n"+
			config.IndexEnd+"\n\n"+
			"## [2026-02-26-000001] Use auto-injection\n\n"+
			"Full decision body here.\n")
	writeTestFile(t, ctxDir, config.FileLearning,
		"# Learnings\n\n"+
			config.IndexStart+"\n"+
			"| Date | Learning |\n"+
			"|------|----------|\n"+
			"| 2026-02-26 | Hooks need JSON output |\n"+
			config.IndexEnd+"\n\n"+
			"## [2026-02-26-000001] Hooks need JSON output\n\n"+
			"Full learning body here.\n")
	writeTestFile(t, ctxDir, config.FileTask,
		"# Tasks\n\n- [ ] Implement context-load-gate v2\n")
	writeTestFile(t, ctxDir, config.FileGlossary,
		"# Glossary\n\n| Term | Definition |\n")

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-inject"}`)
	if err := runContextLoadGate(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)

	// Must be JSON HookResponse
	if !strings.Contains(out, "additionalContext") {
		t.Errorf("expected JSON HookResponse, got: %s", out)
	}

	// Verbatim files: content must appear
	if !strings.Contains(out, "Never commit secrets.") {
		t.Error("expected CONSTITUTION content in output")
	}
	if !strings.Contains(out, "Use filepath.Join.") {
		t.Error("expected CONVENTIONS content in output")
	}
	if !strings.Contains(out, "CLI tool with .context/ directory.") {
		t.Error("expected ARCHITECTURE content in output")
	}
	if !strings.Contains(out, "Work → Reflect → Persist.") {
		t.Error("expected AGENT_PLAYBOOK content in output")
	}

	// Index-only files: index table must appear, full body must NOT
	if !strings.Contains(out, "Use auto-injection") {
		t.Error("expected DECISIONS index entry in output")
	}
	if strings.Contains(out, "Full decision body here.") {
		t.Error("DECISIONS full body should NOT appear — index only")
	}
	if !strings.Contains(out, "Hooks need JSON output") {
		t.Error("expected LEARNINGS index entry in output")
	}
	if strings.Contains(out, "Full learning body here.") {
		t.Error("LEARNINGS full body should NOT appear — index only")
	}

	// Index-only files: must have the "read full entries" label
	if !strings.Contains(out, "index — read full entries by date") {
		t.Error("expected index-only label for DECISIONS/LEARNINGS")
	}

	// TASKS: must NOT be in injected content, but mentioned in footer
	if strings.Contains(out, "Implement context-load-gate v2") {
		t.Error("TASKS content should NOT appear — mention only")
	}
	if !strings.Contains(out, "TASKS.md contains the project") {
		t.Error("expected TASKS mention in footer")
	}

	// GLOSSARY: must NOT appear at all
	if strings.Contains(out, "Glossary") {
		t.Error("GLOSSARY should NOT appear in output")
	}

	// Header and footer
	if !strings.Contains(out, "PROJECT CONTEXT (auto-loaded by system hook") {
		t.Error("expected header in output")
	}
	if !strings.Contains(out, "files loaded") {
		t.Error("expected token count in footer")
	}
}

func TestContextLoadGate_SecondToolUse_Silent(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	// First tool use — injects content and creates marker
	cmd1 := newTestCmd()
	stdin1 := createTempStdin(t, `{"session_id":"test-second-tool"}`)
	if err := runContextLoadGate(cmd1, stdin1); err != nil {
		t.Fatalf("first tool use: unexpected error: %v", err)
	}
	if cmdOutput(cmd1) == "" {
		t.Fatal("first tool use: expected injection output")
	}

	// Second tool use — marker exists, should be silent
	cmd2 := newTestCmd()
	stdin2 := createTempStdin(t, `{"session_id":"test-second-tool"}`)
	if err := runContextLoadGate(cmd2, stdin2); err != nil {
		t.Fatalf("second tool use: unexpected error: %v", err)
	}

	out := cmdOutput(cmd2)
	if out != "" {
		t.Errorf("expected silence on second tool use, got: %s", out)
	}
}

func TestContextLoadGate_DifferentSessions(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	// Session A — injects content
	cmdA := newTestCmd()
	stdinA := createTempStdin(t, `{"session_id":"session-a"}`)
	if err := runContextLoadGate(cmdA, stdinA); err != nil {
		t.Fatalf("session-a: unexpected error: %v", err)
	}
	if cmdOutput(cmdA) == "" {
		t.Fatal("session-a: expected injection output")
	}

	// Verify marker exists for session-a
	marker := filepath.Join(tmpDir, "ctx", "ctx-loaded-session-a")
	if _, err := os.Stat(marker); err != nil {
		t.Errorf("expected marker for session-a, got error: %v", err)
	}

	// Session B — different session_id, should also inject
	cmdB := newTestCmd()
	stdinB := createTempStdin(t, `{"session_id":"session-b"}`)
	if err := runContextLoadGate(cmdB, stdinB); err != nil {
		t.Fatalf("session-b: unexpected error: %v", err)
	}
	outB := cmdOutput(cmdB)
	if !strings.Contains(outB, "PROJECT CONTEXT") {
		t.Errorf("session-b: expected injection, got: %s", outB)
	}

	// Session A again — should be silent
	cmdA2 := newTestCmd()
	stdinA2 := createTempStdin(t, `{"session_id":"session-a"}`)
	if err := runContextLoadGate(cmdA2, stdinA2); err != nil {
		t.Fatalf("session-a repeat: unexpected error: %v", err)
	}
	if cmdOutput(cmdA2) != "" {
		t.Errorf("session-a repeat: expected silence, got: %s", cmdOutput(cmdA2))
	}
}

func TestContextLoadGate_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	// Only write CONSTITUTION — all other files missing
	ctxDir := filepath.Join(workDir, config.DirContext)
	writeTestFile(t, ctxDir, config.FileConstitution,
		"# Constitution\n\nNever commit secrets.\n")

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-missing"}`)
	if err := runContextLoadGate(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	// Should still inject what it can
	if !strings.Contains(out, "Never commit secrets.") {
		t.Error("expected CONSTITUTION content despite other files missing")
	}
	if !strings.Contains(out, "files loaded") {
		t.Error("expected footer with file count")
	}
}

func TestContextLoadGate_NoIndexMarkers(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)

	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	setupContextDir(t)

	ctxDir := filepath.Join(workDir, config.DirContext)
	// DECISIONS without index markers
	writeTestFile(t, ctxDir, config.FileDecision,
		"# Decisions\n\n## [2026-02-26] Some decision\n\nBody.\n")

	cmd := newTestCmd()
	stdin := createTempStdin(t, `{"session_id":"test-no-index"}`)
	if err := runContextLoadGate(cmd, stdin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "(no index entries)") {
		t.Error("expected '(no index entries)' fallback for DECISIONS without markers")
	}
}

func TestExtractIndex(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name: "normal index",
			content: "# Decisions\n\n" +
				config.IndexStart + "\n" +
				"| Date | Decision |\n" +
				"|------|----------|\n" +
				"| 2026-02-26 | Test |\n" +
				config.IndexEnd + "\n",
			want: "| Date | Decision |\n|------|----------|\n| 2026-02-26 | Test |",
		},
		{
			name:    "no markers",
			content: "# Decisions\n\nNo index here.\n",
			want:    "",
		},
		{
			name: "empty index",
			content: "# Decisions\n\n" +
				config.IndexStart + "\n" +
				config.IndexEnd + "\n",
			want: "",
		},
		{
			name: "only start marker",
			content: "# Decisions\n\n" +
				config.IndexStart + "\n" +
				"| orphaned content |\n",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractIndex(tt.content)
			if got != tt.want {
				t.Errorf("extractIndex() = %q, want %q", got, tt.want)
			}
		})
	}
}

// writeTestFile creates a file in dir with the given name and content.
func writeTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600)
	if err != nil {
		t.Fatalf("writeTestFile(%s): %v", name, err)
	}
}
