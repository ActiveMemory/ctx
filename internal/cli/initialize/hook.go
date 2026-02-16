//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/claude"
	"github.com/ActiveMemory/ctx/internal/config"
)

// mergeSettingsPermissions merges ctx permissions into settings.local.json.
//
// Only adds missing permissions to preserve user customizations. Does not
// manage hooks — hook configuration is now provided by the ctx Claude Code
// plugin.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if JSON parsing or file operations fail
func mergeSettingsPermissions(cmd *cobra.Command) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Check if settings.local.json exists
	var settings claude.Settings
	existingContent, err := os.ReadFile(config.FileSettings)
	fileExists := err == nil

	if fileExists {
		if err := json.Unmarshal(existingContent, &settings); err != nil {
			return fmt.Errorf(
				"failed to parse existing %s: %w", config.FileSettings, err,
			)
		}
	}

	// Merge permissions - always additive, never removes existing permissions
	defaultPerms := config.DefaultClaudePermissions
	permsModified := mergePermissions(&settings.Permissions, defaultPerms)

	if !permsModified {
		cmd.Printf(
			"  %s %s (no changes needed)\n", yellow("○"), config.FileSettings,
		)
		return nil
	}

	// Create .claude/ directory if needed
	if err := os.MkdirAll(config.DirClaude, config.PermExec); err != nil {
		return fmt.Errorf("failed to create %s: %w", config.DirClaude, err)
	}

	// Write settings with pretty formatting
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(settings); err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(config.FileSettings, buf.Bytes(), config.PermFile); err != nil {
		return fmt.Errorf("failed to write %s: %w", config.FileSettings, err)
	}

	if fileExists {
		cmd.Printf("  %s %s (added ctx permissions)\n", green("✓"), config.FileSettings)
	} else {
		cmd.Printf("  %s %s\n", green("✓"), config.FileSettings)
	}

	return nil
}

// mergePermissions adds missing permissions to the allow list.
//
// Only adds permissions that don't already exist. Never removes existing
// permissions to preserve user customizations.
//
// Parameters:
//   - perms: Existing permissions config to modify
//   - defaults: Default permissions to add if missing
//
// Returns:
//   - bool: True if any permissions were added
func mergePermissions(perms *claude.PermissionsConfig, defaults []string) bool {
	// Build a set of existing permissions for fast lookup
	existing := make(map[string]bool)
	for _, p := range perms.Allow {
		existing[p] = true
	}

	// Add missing permissions
	added := false
	for _, p := range defaults {
		if !existing[p] {
			perms.Allow = append(perms.Allow, p)
			added = true
		}
	}

	return added
}
