//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package setup

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	coreCopilot "github.com/ActiveMemory/ctx/internal/cli/setup/core/copilot"
	"github.com/spf13/cobra"
)

// TestSetupCommand tests the setup command.
func TestSetupCommand(t *testing.T) {
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
			setupCmd := Cmd()
			setupCmd.SetArgs([]string{tt.tool})

			if err := setupCmd.Execute(); err != nil {
				t.Fatalf("setup %s command failed: %v", tt.tool, err)
			}
		})
	}
}

// TestSetupCommandUnknownTool tests setup command with unknown tool.
func TestSetupCommandUnknownTool(t *testing.T) {
	setupCmd := Cmd()
	setupCmd.SetArgs([]string{"unknown-tool"})

	err := setupCmd.Execute()
	if err == nil {
		t.Error("setup command should fail for unknown tool")
	}
}

// newHookTestCmd creates a cobra command with a captured output buffer.
func newHookTestCmd() *cobra.Command {
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

// setupCmdOutput returns the captured output from a test command.
func setupCmdOutput(cmd *cobra.Command) string {
	return cmd.OutOrStdout().(*bytes.Buffer).String()
}

// TestDeployInstructions_NewFile creates the file from scratch.
func TestDeployInstructions_NewFile(t *testing.T) {
	tmpDir := t.TempDir()

	// coreCopilot.DeployInstructions uses relative paths, so chdir.
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
	if runErr := coreCopilot.DeployInstructions(cmd); runErr != nil {
		t.Fatalf("coreCopilot.DeployInstructions failed: %v", runErr)
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

// TestDeployInstructions_ExistingWithMarker skips when marker exists.
func TestDeployInstructions_ExistingWithMarker(t *testing.T) {
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
	if runErr := coreCopilot.DeployInstructions(cmd); runErr != nil {
		t.Fatalf("coreCopilot.DeployInstructions failed: %v", runErr)
	}

	// File should be unchanged (skipped).
	data, readErr := os.ReadFile(targetFile)
	if readErr != nil {
		t.Fatalf("unexpected read error: %v", readErr)
	}
	if string(data) != existingContent {
		t.Error("file with existing ctx marker should not be modified")
	}

	out := setupCmdOutput(cmd)
	if !strings.Contains(out, "skipped") {
		t.Errorf("output should mention skipped, got: %s", out)
	}
}

// TestDeployInstructions_ExistingWithoutMarker merges content.
func TestDeployInstructions_ExistingWithoutMarker(t *testing.T) {
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
	if runErr := coreCopilot.DeployInstructions(cmd); runErr != nil {
		t.Fatalf("coreCopilot.DeployInstructions failed: %v", runErr)
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

	out := setupCmdOutput(cmd)
	if !strings.Contains(out, "merged") {
		t.Errorf("output should mention merged, got: %s", out)
	}
}
