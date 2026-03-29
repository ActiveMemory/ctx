//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	hookRoot "github.com/ActiveMemory/ctx/internal/cli/hook/cmd/root"
	"github.com/spf13/cobra"
)

// TestHookCommand tests the hook command.
func TestHookCommand(t *testing.T) {
	tests := []struct {
		tool     string
		contains string
	}{
		{"claude-code", "Claude Code Integration"},
		{"cursor", "Cursor IDE Integration"},
		{"aider", "Aider Integration"},
		{"copilot", "GitHub Copilot Integration"},
		{"windsurf", "Windsurf Integration"},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			hookCmd := Cmd()
			hookCmd.SetArgs([]string{tt.tool})

			if err := hookCmd.Execute(); err != nil {
				t.Fatalf("hook %s command failed: %v", tt.tool, err)
			}
		})
	}
}

// TestHookCommandUnknownTool tests hook command with unknown tool.
func TestHookCommandUnknownTool(t *testing.T) {
	hookCmd := Cmd()
	hookCmd.SetArgs([]string{"unknown-tool"})

	err := hookCmd.Execute()
	if err == nil {
		t.Error("hook command should fail for unknown tool")
	}
}

// newHookTestCmd creates a cobra command with a captured output buffer.
func newHookTestCmd() *cobra.Command {
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

// hookCmdOutput returns the captured output from a test command.
func hookCmdOutput(cmd *cobra.Command) string {
	return cmd.OutOrStdout().(*bytes.Buffer).String()
}

// TestWriteCopilotInstructions_NewFile creates the file from scratch.
func TestWriteCopilotInstructions_NewFile(t *testing.T) {
	tmpDir := t.TempDir()

	// hookRoot.WriteCopilotInstructions uses relative paths, so chdir.
	origDir, wdErr := os.Getwd()
	if wdErr != nil {
		t.Fatal(wdErr)
	}
	if chdirErr := os.Chdir(tmpDir); chdirErr != nil {
		t.Fatal(chdirErr)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	cmd := newHookTestCmd()
	if runErr := hookRoot.WriteCopilotInstructions(cmd); runErr != nil {
		t.Fatalf("hookRoot.WriteCopilotInstructions failed: %v", runErr)
	}

	targetFile := filepath.Join(tmpDir, ".github", "copilot-instructions.md")
	data, readErr := os.ReadFile(targetFile)
	if readErr != nil {
		t.Fatalf("expected file to be created: %v", readErr)
	}

	content := string(data)
	if !strings.Contains(content, "<!-- ctx:copilot -->") {
		t.Error("created file should contain the ctx marker")
	}
	if !strings.Contains(content, "<!-- ctx:copilot:end -->") {
		t.Error("created file should contain the ctx end marker")
	}
}

// TestWriteCopilotInstructions_ExistingWithMarker skips when marker exists.
func TestWriteCopilotInstructions_ExistingWithMarker(t *testing.T) {
	tmpDir := t.TempDir()

	origDir, wdErr := os.Getwd()
	if wdErr != nil {
		t.Fatal(wdErr)
	}
	if chdirErr := os.Chdir(tmpDir); chdirErr != nil {
		t.Fatal(chdirErr)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	// Pre-create the file with the ctx marker.
	githubDir := filepath.Join(tmpDir, ".github")
	if mkErr := os.MkdirAll(githubDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	existingContent := "# Existing\n<!-- ctx:copilot -->\nSome content\n"
	targetFile := filepath.Join(githubDir, "copilot-instructions.md")
	if writeErr := os.WriteFile(
		targetFile, []byte(existingContent), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newHookTestCmd()
	if runErr := hookRoot.WriteCopilotInstructions(cmd); runErr != nil {
		t.Fatalf("hookRoot.WriteCopilotInstructions failed: %v", runErr)
	}

	// File should be unchanged (skipped).
	data, readErr := os.ReadFile(targetFile)
	if readErr != nil {
		t.Fatalf("unexpected read error: %v", readErr)
	}
	if string(data) != existingContent {
		t.Error("file with existing ctx marker should not be modified")
	}

	out := hookCmdOutput(cmd)
	if !strings.Contains(out, "skipped") {
		t.Errorf("output should mention skipped, got: %s", out)
	}
}

// TestWriteCopilotInstructions_ExistingWithoutMarker merges content.
func TestWriteCopilotInstructions_ExistingWithoutMarker(t *testing.T) {
	tmpDir := t.TempDir()

	origDir, wdErr := os.Getwd()
	if wdErr != nil {
		t.Fatal(wdErr)
	}
	if chdirErr := os.Chdir(tmpDir); chdirErr != nil {
		t.Fatal(chdirErr)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	// Pre-create file without ctx marker.
	githubDir := filepath.Join(tmpDir, ".github")
	if mkErr := os.MkdirAll(githubDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	originalContent := "# My Custom Copilot Rules\n\nDo cool things.\n"
	targetFile := filepath.Join(githubDir, "copilot-instructions.md")
	writeErr := os.WriteFile(
		targetFile, []byte(originalContent), 0o644)
	if writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newHookTestCmd()
	if runErr := hookRoot.WriteCopilotInstructions(cmd); runErr != nil {
		t.Fatalf("hookRoot.WriteCopilotInstructions failed: %v", runErr)
	}

	data, readErr := os.ReadFile(targetFile)
	if readErr != nil {
		t.Fatalf("unexpected read error: %v", readErr)
	}

	content := string(data)
	// Should contain both the original and the ctx content.
	if !strings.Contains(content, "My Custom Copilot Rules") {
		t.Error("merged file should preserve original content")
	}
	if !strings.Contains(content, "<!-- ctx:copilot -->") {
		t.Error("merged file should contain the ctx marker")
	}

	out := hookCmdOutput(cmd)
	if !strings.Contains(out, "merged") {
		t.Errorf("output should mention merged, got: %s", out)
	}
}
