//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
)

// vscodeDirName is the VS Code workspace configuration directory.
const vscodeDirName = ".vscode"

// CreateVSCodeArtifacts generates VS Code-native configuration files
// as the editor-specific counterpart to Claude Code's settings and hooks.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if file creation fails
func CreateVSCodeArtifacts(cmd *cobra.Command) error {
	if err := os.MkdirAll(vscodeDirName, fs.PermExec); err != nil {
		return fmt.Errorf("failed to create %s/: %w", vscodeDirName, err)
	}

	// .vscode/extensions.json — recommend the ctx extension to collaborators
	if err := writeExtensionsJSON(cmd); err != nil {
		cmd.Println(fmt.Sprintf("  ⚠ extensions.json: %v", err))
	}

	// .vscode/tasks.json — register ctx commands as VS Code tasks
	if err := writeTasksJSON(cmd); err != nil {
		cmd.Println(fmt.Sprintf("  ⚠ tasks.json: %v", err))
	}

	// .vscode/mcp.json — register ctx MCP server for Copilot
	if err := writeMCPJSON(cmd); err != nil {
		cmd.Println(fmt.Sprintf("  ⚠ mcp.json: %v", err))
	}

	return nil
}

func writeExtensionsJSON(cmd *cobra.Command) error {
	target := filepath.Join(vscodeDirName, "extensions.json")

	if _, err := os.Stat(target); err == nil {
		// Exists — check if recommendation is already present
		data, readErr := os.ReadFile(filepath.Clean(target)) //nolint:gosec // path built from constants
		if readErr != nil {
			return readErr
		}
		var existing map[string]interface{}
		if json.Unmarshal(data, &existing) == nil {
			if recs, ok := existing["recommendations"].([]interface{}); ok {
				for _, r := range recs {
					if r == "activememory.ctx-context" {
						cmd.Println(fmt.Sprintf("  ○ %s (recommendation exists)", target))
						return nil
					}
				}
			}
		}
		// File exists but doesn't have our recommendation — leave it alone
		cmd.Println(fmt.Sprintf("  ○ %s (exists, add activememory.ctx-context manually)", target))
		return nil
	}

	content := map[string][]string{
		"recommendations": {"activememory.ctx-context"},
	}
	data, _ := json.MarshalIndent(content, "", "  ")
	data = append(data, '\n')

	if err := os.WriteFile(target, data, fs.PermFile); err != nil {
		return err
	}
	cmd.Println(fmt.Sprintf("  ✓ %s", target))
	return nil
}

func writeTasksJSON(cmd *cobra.Command) error {
	target := filepath.Join(vscodeDirName, "tasks.json")

	if _, err := os.Stat(target); err == nil {
		cmd.Println(fmt.Sprintf("  ○ %s (exists, skipped)", target))
		return nil
	}

	tasks := map[string]interface{}{
		"version": "2.0.0",
		"tasks": []map[string]interface{}{
			{
				"label":   "ctx: status",
				"type":    "shell",
				"command": "ctx status",
				"group":   "none",
				"presentation": map[string]interface{}{
					"reveal": "always",
					"panel":  "shared",
				},
				"problemMatcher": []interface{}{},
			},
			{
				"label":   "ctx: drift",
				"type":    "shell",
				"command": "ctx drift",
				"group":   "none",
				"presentation": map[string]interface{}{
					"reveal": "always",
					"panel":  "shared",
				},
				"problemMatcher": []interface{}{},
			},
			{
				"label":   "ctx: agent",
				"type":    "shell",
				"command": "ctx agent --budget 4000",
				"group":   "none",
				"presentation": map[string]interface{}{
					"reveal": "always",
					"panel":  "shared",
				},
				"problemMatcher": []interface{}{},
			},
			{
				"label":   "ctx: journal",
				"type":    "shell",
				"command": "ctx recall export --all && ctx journal site --build",
				"group":   "none",
				"presentation": map[string]interface{}{
					"reveal": "always",
					"panel":  "shared",
				},
				"problemMatcher": []interface{}{},
			},
			{
				"label":   "ctx: journal-serve",
				"type":    "shell",
				"command": "ctx journal site --serve",
				"group":   "none",
				"presentation": map[string]interface{}{
					"reveal": "always",
					"panel":  "shared",
				},
				"problemMatcher": []interface{}{},
			},
		},
	}
	data, _ := json.MarshalIndent(tasks, "", "  ")
	data = append(data, '\n')

	if err := os.WriteFile(target, data, fs.PermFile); err != nil {
		return err
	}
	cmd.Println(fmt.Sprintf("  ✓ %s", target))
	return nil
}

func writeMCPJSON(cmd *cobra.Command) error {
	target := filepath.Join(vscodeDirName, "mcp.json")

	if _, err := os.Stat(target); err == nil {
		cmd.Println(fmt.Sprintf("  ○ %s (exists, skipped)", target))
		return nil
	}

	mcp := map[string]interface{}{
		"servers": map[string]interface{}{
			"ctx": map[string]interface{}{
				"command": "ctx",
				"args":    []string{"mcp", "serve"},
			},
		},
	}
	data, _ := json.MarshalIndent(mcp, "", "  ")
	data = append(data, '\n')

	if err := os.WriteFile(target, data, fs.PermFile); err != nil {
		return err
	}
	cmd.Println(fmt.Sprintf("  ✓ %s", target))
	return nil
}
