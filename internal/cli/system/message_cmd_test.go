//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/spf13/cobra"
)

// testMessageCmd creates a cobra command for testing message subcommands.
// It captures output in a buffer for assertion.
func testMessageCmd(t *testing.T) (*cobra.Command, *bytes.Buffer) {
	t.Helper()
	cmd := messageCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	return cmd, buf
}

// setupTestDir creates a temp dir, chdir into it, resets rc, and returns a cleanup func.
func setupTestDir(t *testing.T) func() {
	t.Helper()
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	rc.Reset()
	return func() { _ = os.Chdir(origDir) }
}

// --- list tests ---

func TestMessageList_ShowsAllEntries(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"list"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	output := buf.String()
	registry := messages.Registry()
	for _, info := range registry {
		if !strings.Contains(output, info.Hook) {
			t.Errorf("expected hook %q in list output", info.Hook)
		}
		if !strings.Contains(output, info.Variant) {
			t.Errorf("expected variant %q in list output", info.Variant)
		}
	}
}

func TestMessageList_DetectsOverride(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	// Create an override file
	overrideDir := filepath.Join(rc.ContextDir(), "hooks", "messages", "qa-reminder")
	if mkErr := os.MkdirAll(overrideDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	if writeErr := os.WriteFile(filepath.Join(overrideDir, "gate.txt"), []byte("custom"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"list"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "override") {
		t.Error("expected 'override' status in list output when override file exists")
	}
}

func TestMessageList_JSON(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"list", "--json"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	var entries []messageListEntry
	if decodeErr := json.Unmarshal(buf.Bytes(), &entries); decodeErr != nil {
		t.Fatalf("invalid JSON output: %v", decodeErr)
	}

	registry := messages.Registry()
	if len(entries) != len(registry) {
		t.Errorf("expected %d JSON entries, got %d", len(registry), len(entries))
	}

	// Verify first entry has expected fields
	if len(entries) > 0 {
		if entries[0].Hook == "" {
			t.Error("expected non-empty hook in JSON entry")
		}
		if entries[0].Category == "" {
			t.Error("expected non-empty category in JSON entry")
		}
	}
}

func TestMessageList_Categories(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"list"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "customizable") {
		t.Error("expected 'customizable' category in output")
	}
	if !strings.Contains(output, "ctx-specific") {
		t.Error("expected 'ctx-specific' category in output")
	}
}

// --- show tests ---

func TestMessageShow_EmbeddedDefault(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"show", "qa-reminder", "gate"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "Source: embedded default") {
		t.Error("expected 'Source: embedded default' in show output")
	}
	if !strings.Contains(output, "HARD GATE") {
		t.Error("expected embedded template content in show output")
	}
}

func TestMessageShow_UserOverride(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	// Create override
	overrideDir := filepath.Join(rc.ContextDir(), "hooks", "messages", "qa-reminder")
	if mkErr := os.MkdirAll(overrideDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	if writeErr := os.WriteFile(filepath.Join(overrideDir, "gate.txt"), []byte("pytest -x"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"show", "qa-reminder", "gate"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "Source: user override") {
		t.Error("expected 'Source: user override' in show output")
	}
	if !strings.Contains(output, "pytest -x") {
		t.Error("expected override content in show output")
	}
}

func TestMessageShow_TemplateVarsDisplay(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"show", "check-persistence", "nudge"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "{{.PromptsSinceNudge}}") {
		t.Error("expected template variable display in show output")
	}
}

func TestMessageShow_InvalidHook(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cmd, _ := testMessageCmd(t)
	cmd.SetArgs([]string{"show", "nonexistent", "gate"})
	runErr := cmd.Execute()
	if runErr == nil {
		t.Fatal("expected error for invalid hook")
	}
	if !strings.Contains(runErr.Error(), "unknown hook") {
		t.Errorf("expected 'unknown hook' error, got: %v", runErr)
	}
}

func TestMessageShow_InvalidVariant(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cmd, _ := testMessageCmd(t)
	cmd.SetArgs([]string{"show", "qa-reminder", "nonexistent"})
	runErr := cmd.Execute()
	if runErr == nil {
		t.Fatal("expected error for invalid variant")
	}
	if !strings.Contains(runErr.Error(), "unknown variant") {
		t.Errorf("expected 'unknown variant' error, got: %v", runErr)
	}
}

// --- edit tests ---

func TestMessageEdit_CreatesOverride(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"edit", "qa-reminder", "gate"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "Override created") {
		t.Error("expected 'Override created' message")
	}

	// Verify the file was created
	oPath := filepath.Join(rc.ContextDir(), "hooks", "messages", "qa-reminder", "gate.txt")
	data, readErr := os.ReadFile(oPath)
	if readErr != nil {
		t.Fatalf("override file not created: %v", readErr)
	}
	if !strings.Contains(string(data), "HARD GATE") {
		t.Error("override file should contain embedded default content")
	}
}

func TestMessageEdit_RefusesOverwrite(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	// Create existing override
	overrideDir := filepath.Join(rc.ContextDir(), "hooks", "messages", "qa-reminder")
	if mkErr := os.MkdirAll(overrideDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	if writeErr := os.WriteFile(filepath.Join(overrideDir, "gate.txt"), []byte("existing"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd, _ := testMessageCmd(t)
	cmd.SetArgs([]string{"edit", "qa-reminder", "gate"})
	runErr := cmd.Execute()
	if runErr == nil {
		t.Fatal("expected error when override already exists")
	}
	if !strings.Contains(runErr.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", runErr)
	}
}

func TestMessageEdit_CtxSpecificWarning(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"edit", "block-non-path-ctx", "dot-slash"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "ctx-specific") {
		t.Error("expected ctx-specific warning in edit output")
	}
}

func TestMessageEdit_TemplateVarsDisplay(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"edit", "check-persistence", "nudge"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "{{.PromptsSinceNudge}}") {
		t.Error("expected template variable info in edit output")
	}
}

// --- reset tests ---

func TestMessageReset_DeletesOverride(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	// Create override first
	overrideDir := filepath.Join(rc.ContextDir(), "hooks", "messages", "qa-reminder")
	if mkErr := os.MkdirAll(overrideDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	oPath := filepath.Join(overrideDir, "gate.txt")
	if writeErr := os.WriteFile(oPath, []byte("custom"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"reset", "qa-reminder", "gate"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "Override removed") {
		t.Error("expected 'Override removed' message")
	}

	// Verify file was deleted
	if _, statErr := os.Stat(oPath); !os.IsNotExist(statErr) {
		t.Error("expected override file to be deleted")
	}
}

func TestMessageReset_NoOpWhenNoOverride(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cmd, buf := testMessageCmd(t)
	cmd.SetArgs([]string{"reset", "qa-reminder", "gate"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatal(runErr)
	}

	output := buf.String()
	if !strings.Contains(output, "No override found") {
		t.Error("expected 'No override found' message")
	}
}

// --- registry validation ---

func TestRegistryEntriesHaveEmbeddedFiles(t *testing.T) {
	registry := messages.Registry()
	for _, info := range registry {
		_, readErr := assets.HookMessage(info.Hook, info.Variant+".txt")
		if readErr != nil {
			t.Errorf("registry entry %s/%s has no matching embedded file: %v",
				info.Hook, info.Variant, readErr)
		}
	}
}

func TestRegistryCoversAllEmbeddedFiles(t *testing.T) {
	hooks, listErr := assets.ListHookMessages()
	if listErr != nil {
		t.Fatal(listErr)
	}

	for _, hook := range hooks {
		variants, varErr := assets.ListHookVariants(hook)
		if varErr != nil {
			t.Fatalf("failed to list variants for %s: %v", hook, varErr)
		}
		for _, filename := range variants {
			variant := strings.TrimSuffix(filename, ".txt")
			info := messages.Lookup(hook, variant)
			if info == nil {
				t.Errorf("embedded file %s/%s has no registry entry", hook, filename)
			}
		}
	}
}
